package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash  = fmt.Errorf("invalid hash format")
	ErrIncompatible = fmt.Errorf("incompatible version")
	ErrNotMatch     = fmt.Errorf("password not match")
)

// 生成用于存储的认证串
func GenerateAuth(config *config, password string) (string, error) {
	// 生成盐
	salt := make([]byte, config.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, config.Iterations, config.Memory, config.Parallelism, config.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	auth := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, config.Memory, config.Iterations, config.Parallelism, b64Salt, b64Hash,
	)

	return auth, nil
}

// 校验用户输入的密码与存储的认证串是否匹配
func VerifyPassword(password string, auth string) (bool, error) {
	// 解析存储的认证串
	parts := strings.Split(auth, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, ErrInvalidHash
	}
	if parts[2] != "v="+strconv.FormatUint(argon2.Version, 10) {
		return false, ErrIncompatible
	}

	config := DefaultConfig()
	_, err := fmt.Sscanf(
		parts[3], "m=%d,t=%d,p=%d",
		&config.Memory, &config.Iterations, &config.Parallelism,
	)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	config.SaltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	config.KeyLength = uint32(len(hash))

	// 使用相同的参数对输入密码进行哈希
	comparisonHash := argon2.IDKey([]byte(password), salt, config.Iterations, config.Memory, config.Parallelism, config.KeyLength)

	// Gemini Suggested: 使用 constant-time 比较防止计时攻击
	if subtle.ConstantTimeCompare(hash, comparisonHash) != 1 {
		return false, ErrNotMatch
	}

	return true, nil
}
