package nats

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
	"modelo-mcp/internal/config"
	"modelo-mcp/internal/log"
)

type ExampleRequest struct {
	Message string `json:"message"`
}

type ExampleReply struct {
	Echo     string    `json:"echo"`
	Service  string    `json:"service"`
	Ts       time.Time `json:"ts"`
}

func RegisterExampleHandlers(ctx context.Context, js nats.JetStreamContext, cfg *config.Config, logger *slog.Logger) {
	subject := cfg.SubjectRequest
	durable := cfg.NATSDurable

	_, err := js.Subscribe(subject, func(msg *nats.Msg) {
		var req ExampleRequest
		_ = json.Unmarshal(msg.Data, &req)

		reply := ExampleReply{
			Echo:    req.Message,
			Service: cfg.ServiceName,
			Ts:      time.Now().UTC(),
		}
		b, _ := json.Marshal(reply)
		_ = js.Publish(cfg.SubjectReply, b)
		_ = msg.Ack()
		logger.Info("handled message", "subject", subject)
	}, nats.Durable(durable), nats.ManualAck())
	if err != nil {
		logger.Error("subscribe failed", "error", err)
	}
}
