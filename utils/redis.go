package utils

import (
	"context"
	"time"

	"auth-center/config"
)

// RedisClient Redis客户端
var RedisClient = config.RedisClient

// Set 设置键值对
func Set(key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(context.Background(), key, value, expiration).Err()
}

// Get 获取值
func Get(key string) (string, error) {
	return RedisClient.Get(context.Background(), key).Result()
}

// Del 删除键
func Del(key string) error {
	return RedisClient.Del(context.Background(), key).Err()
}

// Exists 检查键是否存在
func Exists(key string) (bool, error) {
	result, err := RedisClient.Exists(context.Background(), key).Result()
	return result > 0, err
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	return RedisClient.Expire(context.Background(), key, expiration).Err()
}

// SAdd 向集合添加成员
func SAdd(key string, members ...interface{}) error {
	return RedisClient.SAdd(context.Background(), key, members...).Err()
}

// SMembers 获取集合所有成员
func SMembers(key string) ([]string, error) {
	return RedisClient.SMembers(context.Background(), key).Result()
}

// SIsMember 检查成员是否在集合中
func SIsMember(key string, member interface{}) (bool, error) {
	return RedisClient.SIsMember(context.Background(), key, member).Result()
}

// SRem 从集合中移除成员
func SRem(key string, members ...interface{}) error {
	return RedisClient.SRem(context.Background(), key, members...).Err()
}

// HSet 设置哈希字段
func HSet(key string, values ...interface{}) error {
	return RedisClient.HSet(context.Background(), key, values...).Err()
}

// HGet 获取哈希字段值
func HGet(key string, field string) (string, error) {
	return RedisClient.HGet(context.Background(), key, field).Result()
}

// HGetAll 获取哈希所有字段
func HGetAll(key string) (map[string]string, error) {
	return RedisClient.HGetAll(context.Background(), key).Result()
}

// HDel 删除哈希字段
func HDel(key string, fields ...string) error {
	return RedisClient.HDel(context.Background(), key, fields...).Err()
}

// LPush 向列表左侧推入元素
func LPush(key string, values ...interface{}) error {
	return RedisClient.LPush(context.Background(), key, values...).Err()
}

// RPop 从列表右侧弹出元素
func RPop(key string) (string, error) {
	return RedisClient.RPop(context.Background(), key).Result()
}

// LLen 获取列表长度
func LLen(key string) (int64, error) {
	return RedisClient.LLen(context.Background(), key).Result()
}

// 缓存键前缀常量
const (
	TokenBlacklistPrefix = "token:blacklist:"
	UserPermissionPrefix = "user:permission:"
	RolePermissionPrefix = "role:permission:"
	APIPermissionPrefix  = "api:permission:"
	AppConfigPrefix      = "app:config:"
)
