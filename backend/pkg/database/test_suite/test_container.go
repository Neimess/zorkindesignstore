package testsuite

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	t_log "github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	pgPort = "5432/tcp"
)

type testContainer struct {
	container testcontainers.Container
	Host      string
	Port      int
}

type dbConfig struct {
	host     string
	user     string
	password string
	dbName   string
}

func New(host, user, password, dbName string) *dbConfig {
	return &dbConfig{
		host:     host,
		user:     user,
		password: password,
		dbName:   dbName,
	}
}

func StartDBContainer(t *testing.T, dbC *dbConfig) *testContainer {
	ctx := context.Background()
	t_log.SetDefault(log.New(io.Discard, "", log.LstdFlags))

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{pgPort},
		Env: map[string]string{
			"POSTGRES_USER":     dbC.user,
			"POSTGRES_PASSWORD": dbC.password,
			"POSTGRES_DB":       dbC.dbName,
		},
		WaitingFor: wait.ForListeningPort(nat.Port(pgPort)).
			WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, pgPort)
	require.NoError(t, err)

	t.Logf("Postgres container started on %s:%s", host, port.Port())

	dsn := buildDSN(dbC.user, dbC.password, dbC.dbName, host, port.Port())
	t.Logf("Attempting to ping DB: %s", dsn)

	err = waitForPostgresReady(dsn, 60*time.Second)
	require.NoError(t, err, "Postgres is not ready after container start")

	return &testContainer{
		container: container,
		Host:      host,
		Port:      port.Int(),
	}
}

func buildDSN(user, pass, dbName, host, port string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbName)
}

func waitForPostgresReady(dsn string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		db, err := sql.Open("postgres", dsn)
		if err == nil {
			defer func() {
				if closeErr := db.Close(); closeErr != nil {
					log.Printf("Failed to close DB connection: %v", closeErr)
				}
			}()
			if pingErr := db.Ping(); pingErr == nil {
				return nil
			}
		}
		time.Sleep(300 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for Postgres to be ready at %s", dsn)
}

func (tc *testContainer) Terminate(t *testing.T) {
	require.NoError(t, tc.container.Terminate(context.Background()))
}
