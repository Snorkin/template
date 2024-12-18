package trace

import (
	"example/pkg/observer/tracing/git"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

type Jaeger struct {
	URL         string `validate:"required"`
	Password    string `validate:"required"`
	Username    string `validate:"required"`
	ServiceName string `validate:"required"`
}

func InitTracer(cfg Jaeger) (*trace.TracerProvider, *jaeger.Exporter, error) {
	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(cfg.URL),
			jaeger.WithUsername(cfg.Username),
			jaeger.WithPassword(cfg.Password),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(git.GetCommitInfo().String()),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, exp, nil
}
