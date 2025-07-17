package testsuite

import (
	"context"
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
		WaitingFor: wait.ForListeningPort(nat.Port(pgPort)).WithStartupTimeout(30 * time.Second),
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

	return &testContainer{
		container: container,
		Host:      host,
		Port:      port.Int(),
	}
}

func (tc *testContainer) Terminate(t *testing.T) {
	require.NoError(t, tc.container.Terminate(context.Background()))
}
