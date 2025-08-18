package nats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"modelo-mcp/internal/config"
	"modelo-mcp/internal/log"
)

type JS struct {
	Conn *nats.Conn
	JS   nats.JetStreamContext
}

func Connect(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*nats.Conn, nats.JetStreamContext, error) {
	nc, err := nats.Connect(cfg.NATSURL, nats.Name(cfg.ServiceName))
	if err != nil { return nil, nil, err }
	js, err := nc.JetStream()
	if err != nil { return nil, nil, err }

	// Ensure stream (idempotent)
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     cfg.NATSJSStream,
		Subjects: []string{"mcp.modelo.>"},
		MaxAge:  7 * 24 * time.Hour,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		logger.Error("add stream", "error", err)
	}
	return nc, js, nil
}
