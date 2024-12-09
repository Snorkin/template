package main

import (
	"context"
	"example/config"
	"example/internal/server/grpc"
	"example/internal/server/http"
	"example/pkg/logger"
	"example/pkg/storage/postgres"
	"github.com/getsentry/sentry-go"
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
	"time"
)

func main() {
	stdLog.Println("Starting service")

	//init config
	config.LoadConfig()
	stdLog.Println("Config loaded")

	//init logger
	logger.InitLogger()

	// app dependencies
	deps := initDeps()
	defer deps.close()

	// start http server
	httpSrv := http.NewServer(deps.pg)
	err := httpSrv.Run()
	if err != nil {
		logger.Log.Fatalf("Failed to start HTTP server: %s", err)
	}
	defer httpSrv.Shutdown()

	// start grpc server
	grpcSrv := grpc.NewServer(deps.pg)
	err = grpcSrv.Run()
	if err != nil {
		logger.Log.Fatalf("Failed to start GRPC server: %s", err)
	}
	defer grpcSrv.Shutdown()

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	<-exitCh

	logger.Log.Info("App is gracefully stopped")
}

type dependencies struct {
	pg  *sqlx.DB
	tp  *trace.TracerProvider
	exp *jaeger.Exporter
}

func initDeps() *dependencies {
	cfg := config.GetConfig()
	//tracing
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(cfg.Jaeger.URL),
			jaeger.WithUsername(cfg.Jaeger.Username),
			jaeger.WithPassword(cfg.Jaeger.Password),
		),
	)
	if err != nil {
		logger.Log.Fatalf("Cannot create Jaeger exporter: %s", err)
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

	//sentry
	if cfg.Sentry.Enabled {
		err = sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.Sentry.Dsn,
			TracesSampleRate: 0.15,
		})
		if err != nil {
			logger.Log.Fatalf("Sentry init error: %s", err)
		} else {
			logger.Log.Info("Sentry connected")
		}
		defer sentry.Flush(time.Second * 5)
	}

	// postgres
	pgDB, err := postgres.InitPsqlDB()
	if err != nil {
		logger.Log.Fatalf("PostgreSQL init error: %s", err)
	} else {
		logger.Log.Infof("PostgreSQL connected")
	}

	return &dependencies{
		pg:  pgDB,
		tp:  tp,
		exp: exporter,
	}
}

func (d *dependencies) close() {
	if d.pg != nil {
		err := d.pg.Close()
		if err != nil {
			logger.Log.Error(err)
		} else {
			logger.Log.Info("PostgreSQL connection resolved")
		}
	}
	if d.tp != nil {
		err := d.tp.Shutdown(context.Background())
		if err != nil {
			logger.Log.Error(err)
		} else {
			logger.Log.Info("Trace provider resolved")
		}
	}
	if d.exp != nil {
		err := d.exp.Shutdown(context.Background())
		if err != nil {
			logger.Log.Error(err)
		} else {
			logger.Log.Info("Exporter resolved")
		}
	}
	logger.Log.Info("All dependencies resolved")
}
