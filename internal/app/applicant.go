package app

import (
	"context"
	handlerhttp "kulturaGo/events-service/internal/handler/http"
	rt "kulturaGo/events-service/internal/handler/routes"
	"kulturaGo/events-service/internal/kafka"
	lg "kulturaGo/events-service/internal/logger"
	"kulturaGo/events-service/internal/repository/postgres"
	"kulturaGo/events-service/internal/scheduler"
	"kulturaGo/events-service/internal/usecase"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Run() {
	lg.Init()

	dbURL := mustEnv("DATABASE_URL")
	kafkaBrokers := strings.Split(mustEnv("KAFKA_BROKERS"), ",")
	kafkaTopic := mustEnv("KAFKA_TOPIC")

	exportEvery := durationEnv("EXPORT_INTERVAL", 5*time.Second)
	batchSize := intEnv("EXPORT_BATCH_SIZE", 500)

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		lg.Log.Fatalf("db: %v", err)
	}
	defer db.Close()

	prod := kafka.New(kafkaBrokers, kafkaTopic)
	defer prod.Close()

	repo := postgres.NewEventRepository(db)
	cursorRepo := postgres.NewCursorRepository(db)

	handler := handlerhttp.NewEventHandler(repo)
	router := rt.NewRoutes(handler)
	server := &http.Server{Addr: ":8090", Handler: router}

	exp := &usecase.Exporter{
		Repo:       repo,
		Producer:   prod,
		OutTopic:   kafkaTopic,
		Batch:      batchSize,
		Interval:   exportEvery,
		Logger:     lg.Log,
		CursorRepo: cursorRepo,
	}
	deact := &scheduler.Deactivator{
		Repo:     repo,
		Interval: 5 * time.Minute,
		Logger:   lg.Log,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go exp.Run(ctx)
	go deact.Run(ctx)
	go func() {
		lg.Log.Info("HTTP listen :8090")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lg.Log.Fatalf("http: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
	lg.Log.Info("shutdown complete")
}

func mustEnv(k string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	lg.Log.Fatalf("env %s required", k)
	return ""
}

func durationEnv(k string, def time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func intEnv(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
