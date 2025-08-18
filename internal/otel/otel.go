package otel

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"modelo-mcp/internal/config"
)

func SetupOTEL(ctx context.Context, cfg *config.Config) (*sdktrace.TracerProvider, *sdktrace.TracerProvider, error) {
	if cfg.OTELExporterEndpoint == "" {
		return nil, nil, nil
	}
	exp, err := otlptrace.New(ctx, otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(cfg.OTELExporterEndpoint),
		otlptracehttp.WithInsecure(),
	))
	if err != nil {
		log.Printf("otel exporter error: %v", err)
		return nil, nil, err
	}
	res, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(cfg.ServiceName),
	))
	if err != nil { return nil, nil, err }
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return tp, tp, nil
}
