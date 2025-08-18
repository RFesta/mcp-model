package http

import (
	"net/http"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"modelo-mcp/internal/config"
	"modelo-mcp/internal/log"
	"modelo-mcp/internal/version"
)

func Router(cfg *config.Config, logger *slog.Logger, promRegistry *prometheus.Registry) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})
	r.Get("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\"service\": \"" + cfg.ServiceName + "\", \"version\": \"" + version.Version + "\", \"commit\": \"" + version.Commit + "\", \"buildTime\": \"" + version.BuildTime + "\"}"))
	})
	r.Handle("/metrics", promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}))

	// Optional pprof (behind flag)
	if cfg.EnablePprof {
		r.Mount("/debug/pprof", middleware.Profiler())
	}

	return r
}
