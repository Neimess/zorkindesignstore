package config

import (
	"log"
	"os"
	"time"

	"github.com/Neimess/zorkin-store-project/pkg/args"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env:"ENV" end-default:"local"`
	Version    string     `yaml:"version" env:"VERSION" end-default:"1.0.0"`
	AdminCode  string     `yaml:"admin_code" env:"ADMIN_CODE" envRequired:"true"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Storage    Storage    `yaml:"storage"`
}

type HTTPServer struct {
	Address         string        `yaml:"address" env:"HTTP_ADDRESS" env-default:"localhost:8080"`
	MaxHeaderBytes  int           `yaml:"max_header_bytes" env:"HTTP_MAX_HEADER_BYTES" env-default:"1048576"` // 1 MB
	ReadTimeout     time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-default:"10s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"5s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"  env:"HTTP_IDLE_TIMEOUT"   env-default:"60s"`
	EnableCORS      bool          `yaml:"enable_cors" env:"HTTP_ENABLE_CORS" env-default:"true"`
	EnablePProf     bool          `yaml:"enable_pprof" env:"HTTP_ENABLE_PPROF" env-default:"false"`
}

type Storage struct {
	Driver string `yaml:"driver"  env:"STORAGE_DRIVER"  end-default:"sqlite"`
	// --- PostgreSQL ---
	Host     string `yaml:"host"              env:"STORAGE_HOST"`
	Port     int    `yaml:"port"              env:"STORAGE_PORT" end-default:"5432"`
	User     string `yaml:"user"              env:"STORAGE_USER"`
	Password string `yaml:"password"          env:"STORAGE_PASSWORD"`
	DBName   string `yaml:"dbname"            env:"STORAGE_DBNAME"`
	SSLMode  string `yaml:"sslmode"           env:"STORAGE_SSLMODE" end-default:"disable"`
	// --- SQLite ---
	Path        string        `yaml:"path"              env:"STORAGE_PATH"  env-default:"./data/app.db"`
	ForeignKeys bool          `yaml:"foreign_keys"      env:"STORAGE_FK"    env-default:"true"`
	BusyTimeout time.Duration `yaml:"busy_timeout"      env:"STORAGE_BT"    env-default:"5s"`
	// --- общие ---
	MaxOpenConns    int           `yaml:"max_open_conns"    env:"STORAGE_MAX_OPEN" end-default:"10"`
	MaxIdleConns    int           `yaml:"max_idle_conns"    env:"STORAGE_MAX_IDLE" end-default:"5"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"STORAGE_LIFETIME" end-default:"30m"`
	ConnMaxTimeout  time.Duration `yaml:"conn_max_timeout"  env:"STORAGE_TIMEOUT"  end-default:"5s"`
}

func MustLoad(argument *args.Args) *Config {
	configPath := fetchConfigPath(argument)
	if configPath == "" {
		log.Fatalf("configuration file path is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("configuration file does not exist: %v", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read configuration file: %v", err)
	}
	return &cfg
}

// fetchConfigPath fetches the configuration file path from command line or environment variable.
// priority: flag > env > default
func fetchConfigPath(a *args.Args) string {
	switch {
	case a.ConfigPath != "":
		return a.ConfigPath
	case os.Getenv("CONFIG_PATH") != "":
		return os.Getenv("CONFIG_PATH")
	default:
		return "./configs/config.yaml"
	}
}
