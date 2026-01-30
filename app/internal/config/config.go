package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPPort  string
	RedisAddr string
	Version   string

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() *Config {
	return &Config{
		HTTPPort:  ":8080",
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		Version:   getEnv("VERSION", "dev"),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}

// getEnv — вспомогательная функция: читает ENV переменную,
// если её нет — возвращает default value (def)
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}