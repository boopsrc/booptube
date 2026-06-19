package downloader

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"booptube/assets"
	"booptube/config"
	"booptube/video"
)

type Client struct {
	cfg config.Config
}

func New(cfg config.Config) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Config() config.Config {
	return c.cfg
}

func (c *Client) EnsureYtdlp(ctx context.Context) error {
	if len(assets.Ytdlp) == 0 {
		return fmt.Errorf("yt-dlp embutido ausente: execute fetch-ytdlp.ps1 ou make fetch-ytdlp antes do build")
	}
	if c.cfg.YtdlpPath == "" {
		return fmt.Errorf("caminho de cache do yt-dlp indisponivel")
	}

	wantSum := sha256.Sum256(assets.Ytdlp)
	wantHex := hex.EncodeToString(wantSum[:])

	if data, err := os.ReadFile(c.cfg.YtdlpPath); err == nil {
		got := sha256.Sum256(data)
		if hex.EncodeToString(got[:]) == wantHex {
			return nil
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := os.MkdirAll(filepath.Dir(c.cfg.YtdlpPath), 0o755); err != nil {
		return fmt.Errorf("criar cache yt-dlp: %w", err)
	}

	tmp := c.cfg.YtdlpPath + ".tmp"
	if err := os.WriteFile(tmp, assets.Ytdlp, 0o755); err != nil {
		return fmt.Errorf("gravar yt-dlp: %w", err)
	}
	if runtime.GOOS != "windows" {
		if err := os.Chmod(tmp, 0o755); err != nil {
			os.Remove(tmp)
			return fmt.Errorf("chmod yt-dlp: %w", err)
		}
	}
	if err := withRetry(5, func() error { return os.Rename(tmp, c.cfg.YtdlpPath) }); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("instalar yt-dlp: %w", err)
	}
	return nil
}

func (c *Client) Download(ctx context.Context, rawURL string, format video.Format) error {
	if c.cfg.DownloadDir == "" {
		return fmt.Errorf("pasta de destino nao definida")
	}
	if err := ensureWritableDir(c.cfg.DownloadDir); err != nil {
		return err
	}
	if err := c.EnsureYtdlp(ctx); err != nil {
		return err
	}

	out := filepath.Join(c.cfg.DownloadDir, "%(title)s.%(ext)s")
	args := []string{
		"--no-playlist",
		"-o", out,
		rawURL,
	}
	switch format {
	case video.FormatMP4:
		args = append([]string{
			"-f", "bv*+ba/b",
			"--merge-output-format", "mp4",
		}, args...)
	case video.FormatMP3:
		args = append([]string{
			"-x",
			"--audio-format", "mp3",
		}, args...)
	default:
		return fmt.Errorf("formato nao suportado: %s", format)
	}

	cmd := exec.CommandContext(ctx, c.cfg.YtdlpPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("download falhou: %w", err)
	}
	return nil
}

func ensureWritableDir(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("criar pasta destino: %w", err)
	}
	test := filepath.Join(dir, ".booptube-write-test")
	if err := os.WriteFile(test, []byte("ok"), 0o644); err != nil {
		return fmt.Errorf("pasta nao gravavel: %w", err)
	}
	return os.Remove(test)
}

func withRetry(max int, fn func() error) error {
	var err error
	backoff := 50 * time.Millisecond
	for i := 0; i < max; i++ {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(backoff)
		backoff *= 2
	}
	return err
}
