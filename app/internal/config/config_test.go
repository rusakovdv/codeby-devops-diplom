package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	os.Setenv("REDIS_ADDR", "1.2.3.4:6379")
	os.Setenv("VERSION", "1.0.0-test")
	defer os.Unsetenv("REDIS_ADDR")
	defer os.Unsetenv("VERSION")

	cfg := Load()

	if cfg.RedisAddr != "1.2.3.4:6379" {
		t.Errorf("expected redis addr 1.2.3.4:6379, got %s", cfg.RedisAddr)
	}

	if cfg.Version != "1.0.0-test" {
		t.Errorf("expected version 1.0.0-test, got %s", cfg.Version)
	}

	if cfg.ReadTimeout != 5*time.Second {
		t.Errorf("expected default read timeout 5s, got %v", cfg.ReadTimeout)
	}
}