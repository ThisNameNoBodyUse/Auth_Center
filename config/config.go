package config

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"auth-center/models"
	"github.com/redis/go-redis/v9"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config 全局配置结构
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Charset  string
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey        string
	TTL              int64
	RefreshSecretKey string
	RefreshTTL       int64
}

var (
	GlobalConfig *Config
	DB           *gorm.DB
	RedisClient  *redis.Client
)

// InitAll 初始化所有配置
func InitAll() {
	loadConfig()
	initDatabase()
	initRedis()
}

// loadConfig 加载配置文件
func loadConfig() {
	// 尝试读取配置文件，如果不存在则使用默认配置
	cfg, err := ini.Load("./config/app.ini")
	if err != nil {
		log.Printf("读取配置文件失败，使用默认配置: %v", err)
		GlobalConfig = getDefaultConfig()
		return
	}

	GlobalConfig = &Config{
		Server: ServerConfig{
			Port: cfg.Section("server").Key("port").MustString("8080"),
			Mode: cfg.Section("server").Key("mode").MustString("debug"),
		},
		Database: DatabaseConfig{
			Host:     cfg.Section("database").Key("host").MustString("localhost"),
			Port:     cfg.Section("database").Key("port").MustString("3306"),
			User:     cfg.Section("database").Key("user").MustString("root"),
			Password: cfg.Section("database").Key("password").MustString(""),
			Database: cfg.Section("database").Key("database").MustString("auth_center"),
			Charset:  cfg.Section("database").Key("charset").MustString("utf8mb4"),
		},
		Redis: RedisConfig{
			Host:     cfg.Section("redis").Key("host").MustString("localhost"),
			Port:     cfg.Section("redis").Key("port").MustString("6379"),
			Password: cfg.Section("redis").Key("password").MustString(""),
			DB:       cfg.Section("redis").Key("db").MustInt(0),
		},
		JWT: JWTConfig{
			SecretKey:        cfg.Section("jwt").Key("secret_key").MustString("your-secret-key"),
			TTL:              cfg.Section("jwt").Key("ttl").MustInt64(3600),
			RefreshSecretKey: cfg.Section("jwt").Key("refresh_secret_key").MustString("your-refresh-secret-key"),
			RefreshTTL:       cfg.Section("jwt").Key("refresh_ttl").MustInt64(7200),
		},
	}
}

// getDefaultConfig 获取默认配置
func getDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_DATABASE", "auth_center"),
			Charset:  getEnv("DB_CHARSET", "utf8mb4"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			SecretKey:        getEnv("JWT_SECRET_KEY", "your-secret-key"),
			TTL:              getEnvInt64("JWT_TTL", 36000),
			RefreshSecretKey: getEnv("JWT_REFRESH_SECRET_KEY", "your-refresh-secret-key"),
			RefreshTTL:       getEnvInt64("JWT_REFRESH_TTL", 72000),
		},
	}
}

// initDatabase 初始化数据库连接
func initDatabase() {
	dsn := buildDSN()

	// 配置 GORM 日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	// 连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移数据库表
	err = DB.AutoMigrate(
		&models.Application{},
		&models.User{},
		&models.SystemAdmin{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.Token{},
		&models.Provider{},
	)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库连接成功!")
}

// initRedis 初始化Redis连接
func initRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     GlobalConfig.Redis.Host + ":" + GlobalConfig.Redis.Port,
		Password: GlobalConfig.Redis.Password,
		DB:       GlobalConfig.Redis.DB,
	})

	// 测试连接
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis连接失败: %v", err)
	}

	log.Println("Redis连接成功!")
}

// buildDSN 构建数据库连接字符串
func buildDSN() string {
	return GlobalConfig.Database.User + ":" +
		GlobalConfig.Database.Password + "@tcp(" +
		GlobalConfig.Database.Host + ":" +
		GlobalConfig.Database.Port + ")/" +
		GlobalConfig.Database.Database + "?charset=" +
		GlobalConfig.Database.Charset + "&parseTime=True&loc=Local"
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return GlobalConfig
}

// 辅助函数
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
