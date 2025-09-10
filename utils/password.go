package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordConfig 密码哈希配置
type PasswordConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

// DefaultPasswordConfig 默认密码配置
var DefaultPasswordConfig = &PasswordConfig{
	Time:    1,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
}

// HashPassword 哈希密码
func HashPassword(password string) (string, error) {
	// 生成随机盐
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 使用Argon2id算法哈希密码
	hash := argon2.IDKey([]byte(password), salt, DefaultPasswordConfig.Time, DefaultPasswordConfig.Memory, DefaultPasswordConfig.Threads, DefaultPasswordConfig.KeyLen)

	// 编码为base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 返回格式化的哈希字符串
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, DefaultPasswordConfig.Memory, DefaultPasswordConfig.Time, DefaultPasswordConfig.Threads, b64Salt, b64Hash)

	return encodedHash, nil
}

// VerifyPassword 验证密码
func VerifyPassword(password, encodedHash string) (bool, error) {
	// 解析哈希字符串
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	// 检查算法
	if parts[1] != "argon2id" {
		return false, fmt.Errorf("unsupported algorithm")
	}

	// 解析版本
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false, err
	}
	if version != argon2.Version {
		return false, fmt.Errorf("incompatible version")
	}

	// 解析参数
	var memory, time uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return false, err
	}

	// 解码盐和哈希
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// 使用相同参数重新计算哈希
	otherHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(hash)))

	// 使用常量时间比较
	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}
