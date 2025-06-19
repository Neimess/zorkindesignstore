package main

import (
	"log/slog"
	"os"

	"github.com/Neimess/zorkin-store-project/internal/config"
	"github.com/Neimess/zorkin-store-project/pkg/args"
	logger "github.com/Neimess/zorkin-store-project/pkg/log"
)

func main() {
	a, cfg := parseArgsAndConfig()
	log := logger.MustInitLogger(cfg.Env)

	switch {
	case a.Force != 0:
		if err := runForce(cfg, log, a.Force); err != nil {
			log.Error("force failed", slog.Any("err", err))
			os.Exit(1)
		}
	case a.Migrate:
		if err := runMigrations(cfg, log); err != nil {
			log.Error("migrate failed", slog.Any("err", err))
			os.Exit(1)
		}
	default:
		log.Info("nothing to do (use --migrate or --force)")
	}
}

func parseArgsAndConfig() (*args.Args, *config.Config) {
	arg := args.Parse()
	return arg, config.MustLoad(arg)
}
