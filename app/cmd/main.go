package main

import (
	"context"
	"example/config"
	"example/internal/server/grpc"
	"example/internal/server/http"
	"example/pkg/observer/logger"
	trace "example/pkg/observer/tracing"
	"example/pkg/storage/postgres"
	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/exporters/jaeger"

	otelTrace "go.opentelemetry.io/otel/sdk/trace"
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
		logger.Log.Fatalp("Failed to start HTTP server: %s", err)
	}
	defer httpSrv.Shutdown()

	// start grpc server
	grpcSrv := grpc.NewServer(deps.pg)
	err = grpcSrv.Run()
	if err != nil {
		logger.Log.Fatalp("Failed to start GRPC server: %s", err)
	}
	defer grpcSrv.Shutdown()

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	<-exitCh

	logger.Log.Infop("App is gracefully stopped")
}

type dependencies struct {
	pg  *sqlx.DB
	tp  *otelTrace.TracerProvider
	exp *jaeger.Exporter
}

func initDeps() *dependencies {
	cfg := config.GetConfig()
	logger.Log.Infoa("Config", cfg)

	//tracing
	tp, exp, err := trace.InitTracer(trace.Jaeger(cfg.Jaeger))
	if err != nil {
		logger.Log.Fatalp("Failed to init tracer: %s", err)
	}

	//sentry
	if cfg.Sentry.Enabled {
		err = sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.Sentry.Dsn,
			TracesSampleRate: 0.15,
		})
		if err != nil {
			logger.Log.Fatalp("Sentry init error: %s", err)
		} else {
			logger.Log.Infop("Sentry connected")
		}
		defer sentry.Flush(time.Second * 5)
	}

	// postgres
	pgDB, err := postgres.InitPsqlDB()
	if err != nil {
		logger.Log.Fatalp("PostgreSQL init error", "error", err)
	} else {
		logger.Log.Infoa("PostgreSQL connected", pgDB.Stats())
	}

	return &dependencies{
		pg:  pgDB,
		tp:  tp,
		exp: exp,
	}
}

func (d *dependencies) close() {
	if d.pg != nil {
		err := d.pg.Close()
		if err != nil {
			logger.Log.Error(err)
		} else {
			logger.Log.Infop("PostgreSQL connection resolved")
		}
	}
	if d.tp != nil {
		err := d.tp.Shutdown(context.Background())
		if err != nil {
			logger.Log.Error(err)
		} else {
			logger.Log.Infop("Trace provider resolved")
		}
	}
	if d.exp != nil {
		err := d.exp.Shutdown(context.Background())
		if err != nil {
			logger.Log.Error(err)
		} else {
			logger.Log.Infop("Exporter resolved")
		}
	}
	logger.Log.Infop("All dependencies resolved")
}
