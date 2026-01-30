package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"myapp/internal/config"
	httpHandlers "myapp/internal/http"
	"myapp/internal/metrics"
	"myapp/internal/redis"
)

func Run() error {
	cfg := config.Load()

	logger, _ := zap.NewProduction()
	log := logger.Sugar()
	defer logger.Sync()

	metrics.Register()

	rdb := redis.New(cfg.RedisAddr)
	defer rdb.Close()

	handler := &httpHandlers.Handler{
		Redis:   rdb,
		Version: cfg.Version,
		Logger:  log,
	}

	server := &http.Server{
		Addr:         cfg.HTTPPort,
		Handler:      handler.Router(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Infof("starting server on %s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Info("shutting down server")

	return server.Shutdown(shutdownCtx)
}
