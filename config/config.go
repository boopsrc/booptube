package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"booptube/assets"
)

type Config struct {
	DownloadDir string `json:"download_dir"`
	YtdlpPath   string `json:"-"`
	FfmpegDir   string `json:"-"`
}

func Load() (Config, error) {
	cfg := Config{
		YtdlpPath: defaultYtdlpPath(),
		FfmpegDir: defaultFfmpegDir(),
	}
	path, err := configPath()
	if err != nil {
		return cfg, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("read config: %w", err)
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config: %w", err)
	}
	cfg.YtdlpPath = defaultYtdlpPath()
	cfg.FfmpegDir = defaultFfmpegDir()
	return cfg, nil
}

func Save(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("user config dir: %w", err)
	}
	return filepath.Join(dir, "booptube", "config.json"), nil
}

func defaultYtdlpPath() string {
	if dir := installToolsDir(); dir != "" {
		return filepath.Join(dir, assets.YtdlpName)
	}
	dir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "booptube", assets.YtdlpName)
}

func defaultFfmpegDir() string {
	if dir := installToolsDir(); dir != "" {
		return dir
	}
	dir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "booptube", "ffmpeg")
}
