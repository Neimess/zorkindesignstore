package main

import (
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/config"
	"github.com/Neimess/zorkin-store-project/pkg/args"
	logger "github.com/Neimess/zorkin-store-project/pkg/log"
)

func main() {
	a, cfg := parseArgsAndConfig()
	log := logger.MustInitLogger(cfg.Env)

	if err := handleMigrations(a, cfg, log); err != nil {
		log.Error("error while migrate:", slog.String("err", err.Error()))
	}
}

func parseArgsAndConfig() (*args.Args, *config.Config) {
	arg := args.Parse()
	return arg, config.MustLoad(arg)
}
