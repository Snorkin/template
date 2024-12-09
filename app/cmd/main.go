package main

import (
	"context"
	"example/config"
	"example/internal/server/grpc"
	"example/internal/server/http"
	"example/pkg/logger"
	"example/pkg/storage/postgres"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	stdLog "log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stdLog.Println("Starting service")

	//init config
	cfg, err := config.LoadConfig()
	if err != nil {
		stdLog.Fatalf("Failed to load config %s", err.Error())
	}
	stdLog.Println("Config loaded")

	//init logger
	log := logger.NewLogger(cfg)
	log.Info("Logger initialized")
	ctx := context.Background()

	// app dependencies
	deps := initDeps(ctx, cfg, log)
	defer deps.close(log)

	// start http server
	httpSrv := http.NewServer(cfg, log, deps.pg)
	err = httpSrv.Run()
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %s", err)
	}
	defer httpSrv.Shutdown()

	// start grpc server
	grpcSrv := grpc.NewServer(cfg, log, deps.pg)
	err = grpcSrv.Run()
	if err != nil {
		log.Fatalf("Failed to start GRPC server: %s", err)
	}
	defer grpcSrv.Shutdown()

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	<-exitCh

	log.Info("App is gracefully stopped")
}

type dependencies struct {
	pg *sqlx.DB
	tp *trace.TracerProvider
}

func initDeps(ctx context.Context, cfg config.Config, log logger.Logger) *dependencies {
	//tracing
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(cfg.Jaeger.URL),
			jaeger.WithUsername(cfg.Jaeger.Username),
			jaeger.WithPassword(cfg.Jaeger.Password),
		),
	)
	if err != nil {
		log.Fatalf("Cannot create Jaeger exporter: %s", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.Jaeger.ServiceName),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// postgres
	pgDB, err := postgres.InitPsqlDB(ctx, cfg.Postgres)
	if err != nil {
		log.Fatalf("PostgreSQL init error: %s", err)
	} else {
		log.Infof("PostgreSQL connected")
	}

	return &dependencies{
		pg: pgDB,
		tp: tp,
	}
}

func (d *dependencies) close(log logger.Logger) {
	if d.pg != nil {
		err := d.pg.Close()
		if err != nil {
			log.Error(err)
		} else {
			log.Info("PostgreSQL connection resolved")
		}
	}
	if d.tp != nil {
		err := d.tp.Shutdown(context.Background())
		if err != nil {
			log.Error(err)
		} else {
			log.Info("Trace provider resolved")
		}
	}
	log.Info("All dependencies resolved")
}
