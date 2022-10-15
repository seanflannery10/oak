package database

import (
	"testing"
)

func TestNew(t *testing.T) {
	cfg := &Config{
		dsn:                   "test",
		minConns:              30,
		maxConns:              30,
		maxConnLifetime:       "60m",
		maxConnLifetimeJitter: "5s",
		maxConnIdleTime:       "30m",
	}

	_, err := New(*cfg)
	if err != nil {
		t.Fatal(err)
	}
}
