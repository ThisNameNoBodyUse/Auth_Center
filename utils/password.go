package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// DefaultCost bcrypt 默认成本因子（范围 4-31，值越大计算越慢，安全性越高）
const DefaultCost = 10

// HashPassword 使用 bcrypt 哈希密码
// 参数：password 原始明文密码
// 返回：加密后的哈希字符串 / 错误信息
func HashPassword(password string) (string, error) {
	// 1. 生成 bcrypt 哈希（自动生成 16 字节随机盐）
	// GenerateFromPassword 会自动处理盐值和成本因子
	hashBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), // 原始密码字节
		DefaultCost,      // 成本因子（推荐 10-12，平衡安全与性能）
	)
	if err != nil {
		return "", err
	}

	// 2. 直接返回哈希字符串（bcrypt 哈希已包含盐值和成本因子，格式为 "$2a$10$盐值$哈希"）
	return string(hashBytes), nil
}

// VerifyPassword 验证 bcrypt 哈希与明文密码是否匹配
// 参数：password 明文密码 / encodedHash 存储的 bcrypt 哈希
// 返回：是否匹配 / 错误信息
func VerifyPassword(password, encodedHash string) (bool, error) {
	// 1. 调用 CompareHashAndPassword 验证（自动解析哈希中的盐值和成本因子）
	err := bcrypt.CompareHashAndPassword(
		[]byte(encodedHash), // 存储的哈希字节
		[]byte(password),    // 待验证的明文密码字节
	)

	// 2. 处理结果（err 为 nil 表示验证成功）
	if err == bcrypt.ErrMismatchedHashAndPassword {
		// 密码不匹配
		return false, nil
	}
	// 其他错误（如哈希格式无效）
	return err == nil, err
}
