package testsuite

import (
	"context"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/Neimess/zorkin-store-project/internal/app"
	"github.com/Neimess/zorkin-store-project/internal/config"
	"github.com/Neimess/zorkin-store-project/pkg/args"
	"github.com/Neimess/zorkin-store-project/pkg/secret/jwt"
	"github.com/stretchr/testify/require"
)

type TestServer struct {
	Cfg    *config.Config
	App    *app.Application
	Server *httptest.Server
}

func RunTestServer(t *testing.T) *TestServer {
	t.Helper()
	root := projectRoot()

	cfg := config.MustLoad(
		&args.Args{
			ConfigPath: fmt.Sprintf("%s/configs/test.yaml", root),
		},
	)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	tc := StartDBContainer(t, New(
		cfg.Storage.Host,
		cfg.Storage.User,
		cfg.Storage.Password,
		cfg.Storage.DBName,
	))

	cfg.Storage.Host = tc.Host
	cfg.Storage.Port = tc.Port

	application, err := app.NewApplication(&app.Deps{
		Config: cfg,
		Logger: logger,
	})
	require.NoError(t, err, "app initialization failed")

	ts := httptest.NewServer(application.ServerHandler())

	t.Cleanup(func() {
		tc.Terminate(t)
		_ = application.Shutdown(context.Background())
		ts.Close()
	})

	return &TestServer{
		Cfg:    cfg,
		App:    application,
		Server: ts,
	}
}

func projectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("failed to get working directory: " + err.Error())
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			panic("project root not found from working dir: " + dir)
		}
		dir = parent
	}
}

func (ts *TestServer) GenerateTestJWT(t *testing.T, userID int64) string {
	token, err := jwt.NewGenerator(jwt.JWTConfig{
		Secret:    ts.Cfg.JWTConfig.JWTSecret,
		Algorithm: ts.Cfg.JWTConfig.Algorithm,
		Issuer:    ts.Cfg.JWTConfig.Issuer,
		Audience:  ts.Cfg.JWTConfig.Audience,
	}).Generate(strconv.FormatInt(userID, 10))

	require.NoError(t, err)
	return token
}
