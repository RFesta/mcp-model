package metrics

import (
	"modelo-mcp/internal/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func NewRegistry(cfg *config.Config) *prometheus.Registry {
	r := prometheus.NewRegistry()
	r.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	r.MustRegister(collectors.NewGoCollector())
	return r
}
