package sqlite

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Option func(*config)

type config struct {
	path         string        // database file path (ignored in memory mode)
	inMemory     bool          // open :memory: database (or named memory DB)
	readOnly     bool          // open in read‑only mode
	foreignKeys  bool          // PRAGMA foreign_keys
	busyTimeout  time.Duration // _busy_timeout pragma
	cache        string        // cache mode (shared|private)
	extraPragmas map[string]string

	// pool settings
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
}

func New(ctx context.Context, opts ...Option) (*sqlx.DB, error) {
	cfg := defaultConfig()
	for _, o := range opts {
		o(cfg)
	}

	dsn := buildDSN(cfg)

	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	db.SetMaxOpenConns(cfg.maxOpenConns)
	db.SetMaxIdleConns(cfg.maxIdleConns)
	db.SetConnMaxLifetime(cfg.connMaxLifetime)

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	return db, nil
}

func WithPath(path string) Option {
	return func(c *config) {
		c.path = path
	}
}

func WithInMemory() Option {
	return func(c *config) {
		c.inMemory = true
		c.path = ":memory:"
	}
}

func WithReadOnly() Option {
	return func(c *config) { c.readOnly = true }
}

func WithForeignKeys(on bool) Option {
	return func(c *config) { c.foreignKeys = on }
}
func WithBusyTimeout(d time.Duration) Option {
	return func(c *config) { c.busyTimeout = d }
}
func WithCache(mode string) Option {
	return func(c *config) { c.cache = mode }
}

func WithPragma(key, value string) Option {
	return func(c *config) { c.extraPragmas[key] = value }
}

func WithPool(maxOpen, maxIdle int, maxLifetime time.Duration) Option {
	return func(c *config) {
		c.maxOpenConns = maxOpen
		c.maxIdleConns = maxIdle
		c.connMaxLifetime = maxLifetime
	}
}

func buildDSN(c *config) string {
	var b strings.Builder

	// base path / identifier
	if c.inMemory {
		b.WriteString("file:")
		if c.path != "" {
			b.WriteString(c.path)
		} else {
			b.WriteString("memdb1")
		}
	} else {
		b.WriteString(c.path)
	}

	params := url.Values{}

	switch {
	case c.inMemory:
		params.Set("mode", "memory")
	case c.readOnly:
		params.Set("mode", "ro")
	}

	if c.cache != "" {
		params.Set("cache", c.cache)
	}

	if c.busyTimeout > 0 {
		params.Set("_busy_timeout", fmt.Sprintf("%d", int(c.busyTimeout.Milliseconds())))
	}

	if c.foreignKeys {
		params.Set("_foreign_keys", "on")
	} else {
		params.Set("_foreign_keys", "off")
	}

	for k, v := range c.extraPragmas {
		params.Set(k, v)
	}

	if q := params.Encode(); q != "" {
		b.WriteString("?")
		b.WriteString(q)
	}

	return b.String()
}

func defaultConfig() *config {
	return &config{
		path:            "file.db",
		foreignKeys:     true,
		busyTimeout:     5 * time.Second,
		cache:           "shared",
		extraPragmas:    map[string]string{},
		maxOpenConns:    1,
		maxIdleConns:    1,
		connMaxLifetime: 0,
	}
}
