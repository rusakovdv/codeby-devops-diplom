package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func setupTestEnv() (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return rdb, mr
}

func TestHandler_Home(t *testing.T) {
	rdb, mr := setupTestEnv()
	defer mr.Close()

	logger := zap.NewNop().Sugar()

	h := &Handler{
		Redis:   rdb,
		Version: "test-v1",
		Logger:  logger,
	}

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	h.Router().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	val, _ := mr.Get("hits")
	if val != "1" {
		t.Errorf("expected hits to be 1, got %s", val)
	}

	body := rr.Body.String()
	expectedVer := "Version: test-v1"

	if !strings.Contains(body, expectedVer) {
		t.Errorf("body does not contain version: got %s", body)
	}
}

func TestHandler_Healthz(t *testing.T) {
	rdb, mr := setupTestEnv()
	defer mr.Close()
	logger := zap.NewNop().Sugar()

	h := &Handler{Redis: rdb, Logger: logger}

	req, _ := http.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	h.Router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("liveness probe failed")
	}
}

func TestHandler_Readyz(t *testing.T) {
	rdb, mr := setupTestEnv()
	defer mr.Close()
	logger := zap.NewNop().Sugar()
	h := &Handler{Redis: rdb, Logger: logger}

	req, _ := http.NewRequest("GET", "/readyz", nil)
	rr := httptest.NewRecorder()
	h.Router().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("readiness probe failed when redis is up")
	}

	mr.Close()
	rr2 := httptest.NewRecorder()
	h.Router().ServeHTTP(rr2, req)
	if rr2.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 when redis is down, got %v", rr2.Code)
	}
}
