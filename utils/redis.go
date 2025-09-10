package utils

import (
	"context"
	"errors"
	"time"

	"auth-center/config"
)

// Set 设置键值对
func Set(key string, value interface{}, expiration time.Duration) error {
	if config.RedisClient == nil {
		return errors.New("redis client is nil (not initialized)")
	}
	return config.RedisClient.Set(context.Background(), key, value, expiration).Err()
}

// Get 获取值
func Get(key string) (string, error) {
	if config.RedisClient == nil {
		return "", errors.New("redis client is nil (not initialized)")
	}
	return config.RedisClient.Get(context.Background(), key).Result()
}

// Del 删除键
func Del(key string) error {
	if config.RedisClient == nil {
		return errors.New("redis client is nil (not initialized)")
	}
	return config.RedisClient.Del(context.Background(), key).Err()
}

// Exists 检查键是否存在
func Exists(key string) (bool, error) {
	if config.RedisClient == nil {
		return false, errors.New("redis client is nil (not initialized)")
	}
	result, err := config.RedisClient.Exists(context.Background(), key).Result()
	return result > 0, err
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	if config.RedisClient == nil {
		return errors.New("redis client is nil (not initialized)")
	}
	return config.RedisClient.Expire(context.Background(), key, expiration).Err()
}

// SAdd 向集合添加成员
func SAdd(key string, members ...interface{}) error {
	if config.RedisClient == nil {
		return errors.New("redis client is nil (not initialized)")
	}
	return config.RedisClient.SAdd(context.Background(), key, members...).Err()
}

// SMembers 获取集合所有成员
func SMembers(key string) ([]string, error) {
	if config.RedisClient == nil {
		return nil, errors.New("redis client is nil (not initialized)")
	}
	return config.RedisClient.SMembers(context.Background(), key).Result()
}

// SIsMember 检查成员是否在集合中
func SIsMember(key string, member interface{}) (bool, error) {
	if config.RedisClient == nil {
		return false, errors.New("redis client is nil (not initialized)")
	}
	return config.RedisClient.SIsMember(context.Background(), key, member).Result()
}

// SRem 从集合中移除成员
func SRem(key string, members ...interface{}) error {
	if config.RedisClient == nil {
		return errors.New("redis client is nil (not initialized)")
	}
	return config.RedisClient.SRem(context.Background(), key, members...).Err()
}

// HSet 设置哈希字段
func HSet(key string, values ...interface{}) error {
	if config.RedisClient == nil {
		return errors.New("redis client is nil (not initialized)")
	}
	return config.RedisClient.HSet(context.Background(), key, values...).Err()
}

// HGet 获取哈希字段值
func HGet(key string, field string) (string, error) {
	if config.RedisClient == nil {
		return "", errors.New("redis client is nil (not initialized)")
	}
	return config.RedisClient.HGet(context.Background(), key, field).Result()
}

// HGetAll 获取哈希所有字段
func HGetAll(key string) (map[string]string, error) {
	if config.RedisClient == nil {
		return nil, errors.New("redis client is nil (not initialized)")
	}
	return config.RedisClient.HGetAll(context.Background(), key).Result()
}

// HDel 删除哈希字段
func HDel(key string, fields ...string) error {
	if config.RedisClient == nil {
		return errors.New("redis client is nil (not initialized)")
	}
	return config.RedisClient.HDel(context.Background(), key, fields...).Err()
}

// LPush 向列表左侧推入元素
func LPush(key string, values ...interface{}) error {
	if config.RedisClient == nil {
		return errors.New("redis client is nil (not initialized)")
	}
	return config.RedisClient.LPush(context.Background(), key, values...).Err()
}

// RPop 从列表右侧弹出元素
func RPop(key string) (string, error) {
	if config.RedisClient == nil {
		return "", errors.New("redis client is nil (not initialized)")
	}
	return config.RedisClient.RPop(context.Background(), key).Result()
}

// LLen 获取列表长度
func LLen(key string) (int64, error) {
	if config.RedisClient == nil {
		return 0, errors.New("redis client is nil (not initialized)")
	}
	return config.RedisClient.LLen(context.Background(), key).Result()
}

// 缓存键前缀常量
const (
	TokenBlacklistPrefix = "token:blacklist:"
	UserPermissionPrefix = "user:permission:"
	RolePermissionPrefix = "role:permission:"
	APIPermissionPrefix  = "api:permission:"
	AppConfigPrefix      = "app:config:"
)
