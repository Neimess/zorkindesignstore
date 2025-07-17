// cmd/migrate/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Neimess/zorkin-store-project/pkg/migrator"
)

func main() {
	var (
		up    = flag.Bool("up", false, "Apply all up migrations")
		down  = flag.Bool("down", false, "Rollback the last migration")
		force = flag.Int("force", 0, "Force to specific version")
		dsn   = flag.String("dsn", os.Getenv("DSN"), "Database DSN (can also be set via DSN env var)")
	)
	flag.Parse()

	if *dsn == "" {
		log.Fatal("Missing DSN. Set via --dsn or DSN env var")
	}

	var mode migrator.Mode
	switch {
	case *up:
		mode = migrator.Up
	case *down:
		mode = migrator.DownOne
	case *force != 0:
		mode = migrator.ForceTo
	default:
		fmt.Println("No action specified. Use -up, -down or -force N")
		return
	}

	if err := migrator.Run(*dsn, migrator.Options{
		Mode:    mode,
		Version: *force,
	}); err != nil {
		log.Fatalf("migrate: %v", err)
	}
}
