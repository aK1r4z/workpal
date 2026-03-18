package auth

import (
	"os"
	"strconv"
)

type config struct {
	Memory      uint32
	Iterations   uint32
	Parallelism uint8

	SaltLength uint32
	KeyLength  uint32
}

func DefaultConfig() *config {
	return &config{
		Memory:      64 * 1024,
		Iterations:   3,
		Parallelism: 2,

		SaltLength: 16,
		KeyLength:  32,
	}
}

func (c *config) Load() {
	if m := os.Getenv("AUTH_ARGON2_MEMORY"); m != "" {
		if val, err := strconv.ParseUint(m, 10, 32); err == nil {
			c.Memory = uint32(val)
		}
	}
	if t := os.Getenv("AUTH_ARGON2_ITERATION"); t != "" {
		if val, err := strconv.ParseUint(t, 10, 32); err == nil {
			c.Iterations = uint32(val)
		}
	}
	if p := os.Getenv("AUTH_ARGON2_MEMORY"); p != "" {
		if val, err := strconv.ParseUint(p, 10, 8); err == nil {
			c.Parallelism = uint8(val)
		}
	}

	// 虽然提供了选项，但不建议修改
	if s := os.Getenv("AUTH_ARGON2_SALT_LENGTH"); s != "" {
		if val, err := strconv.ParseUint(s, 10, 32); err == nil {
			c.SaltLength = uint32(val)
		}
	}
	if k := os.Getenv("AUTH_ARGON2_KEY_LENGTH"); k != "" {
		if val, err := strconv.ParseUint(k, 10, 32); err == nil {
			c.KeyLength = uint32(val)
		}
	}
}

var Config = DefaultConfig()

var Pepper = ""
