package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
	maxLifetime  string
}

const defaultTimeout = 3 * time.Second

func New(cfg Config) (*sql.DB, error) {

	db, err := sql.Open("pgx", "postgres://"+cfg.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.maxOpenConns)
	db.SetMaxIdleConns(cfg.maxIdleConns)
	db.SetMaxOpenConns(cfg.maxOpenConns)
	db.SetMaxIdleConns(cfg.maxIdleConns)

	idleDuration, err := time.ParseDuration(cfg.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(idleDuration)

	lifeDuration, err := time.ParseDuration(cfg.maxLifetime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(lifeDuration)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
