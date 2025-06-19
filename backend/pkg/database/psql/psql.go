package psql

import (
	"context"
	"fmt"
	"net/url"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Option func(*config)

func WithHost(h string) Option      { return func(c *config) { c.host = h } }
func WithPort(p int) Option         { return func(c *config) { c.port = p } }
func WithUser(u string) Option      { return func(c *config) { c.user = u } }
func WithPassword(pw string) Option { return func(c *config) { c.password = pw } }
func WithDB(db string) Option       { return func(c *config) { c.dbname = db } }
func WithSSL(mode string) Option    { return func(c *config) { c.sslmode = mode } }
func WithMaxConns(max int) Option   { return func(c *config) { c.maxOpenConns = max } }
func WithMinConns(min int) Option   { return func(c *config) { c.minIdleConns = min } }
func WithConnLifetime(d time.Duration) Option {
	return func(c *config) { c.connMaxLifetime = d }
}
func WithConnectTimeout(d time.Duration) Option {
	return func(c *config) { c.connectTimeout = d }
}

func New(ctx context.Context, opts ...Option) (*sqlx.DB, error) {
	cfg := defaultConfig()
	for _, option := range opts {
		option(cfg)
	}

	dsn, err := cfg.buildDSN()
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.maxOpenConns)
	db.SetMaxIdleConns(cfg.minIdleConns)
	db.SetConnMaxLifetime(cfg.connMaxLifetime)

	pingCtx, cancel := context.WithTimeout(ctx, cfg.connectTimeout)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return db, nil
}

func (c *config) buildDSN() (string, error) {
	if c.host == "" || c.user == "" || c.dbname == "" {
		return "", fmt.Errorf("incomplete database config: host/user/dbname must be set")
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(c.user),
		url.QueryEscape(c.password),
		c.host,
		c.port,
		c.dbname,
		c.sslmode,
	), nil
}

type config struct {
	host            string
	port            int
	user            string
	password        string
	dbname          string
	sslmode         string
	maxOpenConns    int
	minIdleConns    int
	connMaxLifetime time.Duration
	connectTimeout  time.Duration
}

func defaultConfig() *config {
	return &config{
		host:            "localhost",
		port:            5432,
		user:            "postgres",
		password:        "",
		dbname:          "postgres",
		sslmode:         "disable",
		maxOpenConns:    10,
		minIdleConns:    5,
		connMaxLifetime: 30 * time.Minute,
		connectTimeout:  5 * time.Second,
	}
}
