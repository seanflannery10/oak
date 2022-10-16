package database

import (
	"testing"
)

func TestNew(t *testing.T) {
	_, err := New(Config{
		dsn:                   "postgresql://testing",
		minConns:              30,
		maxConns:              30,
		maxConnLifetime:       "60m",
		maxConnLifetimeJitter: "5s",
		maxConnIdleTime:       "30m",
	})
	if err != nil {
		t.Fatal(err)
	}
}
