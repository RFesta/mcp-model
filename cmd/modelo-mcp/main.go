package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"modelo-mcp/internal/config"
	"modelo-mcp/internal/log"
	"modelo-mcp/internal/metrics"
	natsx "modelo-mcp/internal/nats"
	"modelo-mcp/internal/otel"
	httpx "modelo-mcp/internal/transport/http"
	"modelo-mcp/internal/version"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()

	logger := log.New(cfg.Env, cfg.ServiceName)
	logger.Info("starting service", "service", cfg.ServiceName, "version", version.Version, "commit", version.Commit)

	// OTEL (tracing and metrics)
	tp, mp, err := otel.SetupOTEL(ctx, cfg)
	if err != nil {
		logger.Error("OTEL setup failed", "error", err)
	}
	defer func() {
		if tp != nil { _ = tp.Shutdown(context.Background()) }
		if mp != nil { _ = mp.Shutdown(context.Background()) }
	}()

	// Metrics registry
	promReg := metrics.NewRegistry(cfg)

	// NATS JetStream
	nc, jsm, err := natsx.Connect(ctx, cfg, logger)
	if err != nil {
		logger.Error("failed to connect NATS", "error", err)
		os.Exit(1)
	}
	defer nc.Drain()

	// HTTP server
	r := httpx.Router(cfg, logger, promReg)
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:           otelhttp.NewHandler(r, cfg.ServiceName),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	go func() {
		logger.Info("http server listening", "port", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("http server failed", "error", err)
		}
	}()

	// Handlers (example)
	natsx.RegisterExampleHandlers(ctx, jsm, cfg, logger)

	// Signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	logger.Info("shutdown signal received")

	shutdownCtx, cancel2 := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel2()
	_ = server.Shutdown(shutdownCtx)
	logger.Info("bye")
}
