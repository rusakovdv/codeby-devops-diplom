package http

import (
	"context"
	"fmt"
	nethttp "net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"myapp/internal/metrics"
)

type Handler struct {
	Redis   *redis.Client
	Version string
	Logger  *zap.SugaredLogger
}

func (h *Handler) LoggingMiddleware(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		start := time.Now()

		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: nethttp.StatusOK}

		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start)

		h.Logger.Infow("HTTP Request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrappedWriter.statusCode,
			"duration", duration.String(),
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
	})
}

type responseWriter struct {
	nethttp.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (h *Handler) Router() nethttp.Handler {
	mux := nethttp.NewServeMux()

	mux.HandleFunc("/", h.home)
	mux.HandleFunc("/healthz", h.liveness)
	mux.HandleFunc("/readyz", h.readiness)
	mux.Handle("/metrics", promhttp.Handler())

	return h.LoggingMiddleware(mux)
}

func (h *Handler) home(w nethttp.ResponseWriter, r *nethttp.Request) {
	metrics.Requests.WithLabelValues("/").Inc()

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	count, err := h.Redis.Incr(ctx, "hits").Result()
	if err != nil {
		h.Logger.Errorw("redis error", "error", err)
		nethttp.Error(w, "redis error", 500)
		return
	}

	hostname, _ := os.Hostname()

	fmt.Fprintf(w,
		"Hello from DevOps Diploma!!!\nHostname: %s\nVersion: %s\nHits: %d\n",
		hostname, h.Version, count,
	)
}

func (h *Handler) liveness(w nethttp.ResponseWriter, _ *nethttp.Request) {
	w.WriteHeader(nethttp.StatusOK)
}

func (h *Handler) readiness(w nethttp.ResponseWriter, _ *nethttp.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := h.Redis.Ping(ctx).Err(); err != nil {
		w.WriteHeader(nethttp.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(nethttp.StatusOK)
}
