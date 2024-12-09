package postgres

import (
	"context"
	"example/config"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
)

func InitPsqlDB(ctx context.Context, cfg config.Postgres) (*sqlx.DB, error) {
	_, span := otel.Tracer("").Start(ctx, "storage.InitPsqlDB")
	defer span.End()

	connectionURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	database, err := sqlx.Open(cfg.PGDriver, connectionURL)
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(cfg.Settings.MaxOpenConns)
	database.SetConnMaxLifetime(cfg.Settings.ConnMaxLifetime)
	database.SetMaxIdleConns(cfg.Settings.MaxIdleConns)
	database.SetConnMaxIdleTime(cfg.Settings.ConnMaxIdleTime)

	if err = database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}
