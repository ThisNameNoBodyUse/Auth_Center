package utils

import (
    "crypto/rand"
)

var codeAlphabet = []byte("ABCDEFGHJKLMNPQRSTUVWXYZ23456789abcdefghijkmnpqrstuvwxyz")

// GenerateShortCode 生成较短且碰撞概率低的编码
// length 建议 8-12
func GenerateShortCode(length int) string {
    if length <= 0 {
        length = 10
    }
    b := make([]byte, length)
    // 使用 crypto/rand 填充随机索引
    for i := 0; i < length; i++ {
        // 生成一个安全随机字节，并映射到字母表
        var rb [1]byte
        _, err := rand.Read(rb[:])
        if err != nil {
            // 退化处理：用常量字符，避免崩溃
            b[i] = codeAlphabet[i%len(codeAlphabet)]
            continue
        }
        b[i] = codeAlphabet[int(rb[0])%len(codeAlphabet)]
    }
    return string(b)
}


