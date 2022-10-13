package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Config struct {
	dsn                   string
	maxConns              int32
	maxConnLifetime       string
	maxConnLifetimeJitter string
	maxConnIdleTime       string
}

func New(cfg Config) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.dsn)
	if err != nil {
		return nil, err
	}

	config.MaxConns = cfg.maxConns

	maxConnLifetime, err := time.ParseDuration(cfg.maxConnLifetime)
	if err != nil {
		return nil, err
	}

	config.MaxConnLifetime = maxConnLifetime

	maxConnLifetimeJitter, err := time.ParseDuration(cfg.maxConnLifetimeJitter)
	if err != nil {
		return nil, err
	}

	config.MaxConnLifetimeJitter = maxConnLifetimeJitter

	maxConnIdleTime, err := time.ParseDuration(cfg.maxConnIdleTime)
	if err != nil {
		return nil, err
	}

	config.MaxConnIdleTime = maxConnIdleTime

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	defer pool.Close()

	return pool, nil
}
