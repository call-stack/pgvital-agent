package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	DatabaseURL string
	ServerURL   string
	APIKey      string
	Interval    time.Duration
	Hostname    string
	DataDir     string
	Insecure    bool
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL: os.Getenv("PGVITALS_DATABASE_URL"),
		ServerURL:   os.Getenv("PGVITALS_SERVER_URL"),
		APIKey:      os.Getenv("PGVITALS_API_KEY"),
		Interval:    5 * time.Minute,
		Hostname:    os.Getenv("PGVITALS_HOSTNAME"),
		DataDir:     os.Getenv("PGVITALS_DATA_DIR"),
		Insecure:    os.Getenv("PGVITALS_INSECURE") == "true",
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("PGVITALS_DATABASE_URL is required")
	}
	if cfg.ServerURL == "" {
		cfg.ServerURL = "https://api.pgvitals.kafal.studio"
	}
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("PGVITALS_API_KEY is required")
	}

	if v := os.Getenv("PGVITALS_INTERVAL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("invalid PGVITALS_INTERVAL: %w", err)
		}
		if d < 60*time.Second || d > 3600*time.Second {
			return nil, fmt.Errorf("PGVITALS_INTERVAL must be between 1m and 1h")
		}
		cfg.Interval = d
	}

	if cfg.Hostname == "" {
		h, _ := os.Hostname()
		cfg.Hostname = h
	}

	if cfg.DataDir == "" {
		home, _ := os.UserHomeDir()
		cfg.DataDir = home + "/.pgvitals"
	}

	return cfg, nil
}
