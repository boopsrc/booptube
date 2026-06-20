package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type WebConfig struct {
	Addr           string
	DownloadDir    string
	MaxConcurrent  int
	FileTTL        time.Duration
	DownloadTimeout time.Duration
	YtdlpPath      string
	FfmpegDir      string
}

func LoadWeb() (WebConfig, error) {
	cfg := WebConfig{
		Addr:            envOr("BOOPTUBE_ADDR", ":8080"),
		DownloadDir:     envOr("BOOPTUBE_DOWNLOAD_DIR", "/data/downloads"),
		MaxConcurrent:   envIntOr("BOOPTUBE_MAX_CONCURRENT", 4),
		FileTTL:         envDurationOr("BOOPTUBE_FILE_TTL", 10*time.Minute),
		DownloadTimeout: envDurationOr("BOOPTUBE_DOWNLOAD_TIMEOUT", 30*time.Minute),
		YtdlpPath:       defaultYtdlpPath(),
		FfmpegDir:       defaultFfmpegDir(),
	}
	if cfg.MaxConcurrent < 1 {
		return cfg, fmt.Errorf("BOOPTUBE_MAX_CONCURRENT deve ser >= 1")
	}
	if cfg.FileTTL < time.Minute {
		return cfg, fmt.Errorf("BOOPTUBE_FILE_TTL deve ser >= 1m")
	}
	return cfg, nil
}

func (w WebConfig) ToConfig() Config {
	return Config{
		DownloadDir: w.DownloadDir,
		YtdlpPath:   w.YtdlpPath,
		FfmpegDir:   w.FfmpegDir,
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envIntOr(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func envDurationOr(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
