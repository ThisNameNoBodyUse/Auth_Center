package main

import (
	"log"
	"time"

	"auth-center/config"
	"auth-center/routers"
	"auth-center/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 认证授权中心 API
// @version 1.0
// @description 提供统一的认证授权服务，支持多应用、多租户
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 初始化配置
	config.InitAll()

	// 创建 Gin 实例
	r := gin.Default()

	// 自定义 CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源，生产环境应该限制
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-App-Id", "X-App-Secret"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 处理预检请求
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-App-Id, X-App-Secret")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Status(200)
	})

	// 添加 Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "auth-center"})
	})

	// 初始化路由
	routers.InitRoutes(r)

	// 自动生成 swagger 文档
	err := utils.RunSwagInit()
	if err != nil {
		log.Printf("Swagger 文档生成失败: %v\n", err)
	}

	// 启动服务
	port := config.GetConfig().Server.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("认证授权中心启动在端口: %s", port)
	r.Run(":" + port)
}
