package main

import (
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Neimess/zorkin-store-project/migrations"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
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
		log.Fatal("Missing DSN. Set via --dsn or DSN environment variable")
	}

	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("create postgres driver: %v", err)
	}

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatalf("create iofs source: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		log.Fatalf("init migrate: %v", err)
	}

	switch {
	case *force != 0:
		if err := m.Force(*force); err != nil {
			log.Fatalf("force: %v", err)
		}
		fmt.Printf("✔ Forced to version %d\n", *force)

	case *down:
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("down: %v", err)
		}
		fmt.Println("✔ Rolled back one step")

	case *up:
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("up: %v", err)
		}
		version, _, _ := m.Version()
		fmt.Printf("✔ Migrated to version %d\n", version)

	default:
		fmt.Println("No action specified. Use -up, -down or -force N")
	}
}
