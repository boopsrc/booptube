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
		return fmt.Errorf("yt-dlp embutido ausente: execute scripts/fetch-ytdlp.ps1, scripts/fetch-ytdlp.sh ou make fetch-ytdlp antes do build")
	}
	if c.cfg.YtdlpPath == "" {
		return fmt.Errorf("caminho de cache do yt-dlp indisponivel")
	}
	return ensureBinary(ctx, assets.Ytdlp, c.cfg.YtdlpPath)
}

func (c *Client) EnsureFfmpeg(ctx context.Context) error {
	if len(assets.Ffmpeg) == 0 || len(assets.Ffprobe) == 0 {
		return fmt.Errorf("ffmpeg embutido ausente: execute scripts/fetch-ffmpeg.ps1, scripts/fetch-ffmpeg.sh ou make fetch-ffmpeg antes do build")
	}
	if c.cfg.FfmpegDir == "" {
		return fmt.Errorf("caminho de cache do ffmpeg indisponivel")
	}
	if err := os.MkdirAll(c.cfg.FfmpegDir, 0o755); err != nil {
		return fmt.Errorf("criar cache ffmpeg: %w", err)
	}
	ffmpegPath := filepath.Join(c.cfg.FfmpegDir, assets.FfmpegName)
	ffprobePath := filepath.Join(c.cfg.FfmpegDir, assets.FfprobeName)
	if err := ensureBinary(ctx, assets.Ffmpeg, ffmpegPath); err != nil {
		return fmt.Errorf("instalar ffmpeg: %w", err)
	}
	if err := ensureBinary(ctx, assets.Ffprobe, ffprobePath); err != nil {
		return fmt.Errorf("instalar ffprobe: %w", err)
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
	if err := c.EnsureFfmpeg(ctx); err != nil {
		return err
	}

	out := filepath.Join(c.cfg.DownloadDir, "%(title)s.%(ext)s")
	args := []string{
		"--no-playlist",
		"--ffmpeg-location", c.cfg.FfmpegDir,
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

func ensureBinary(ctx context.Context, data []byte, dest string) error {
	wantSum := sha256.Sum256(data)
	wantHex := hex.EncodeToString(wantSum[:])

	if existing, err := os.ReadFile(dest); err == nil {
		got := sha256.Sum256(existing)
		if hex.EncodeToString(got[:]) == wantHex {
			return nil
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("criar diretorio: %w", err)
	}

	tmp := dest + ".tmp"
	if err := os.WriteFile(tmp, data, 0o755); err != nil {
		return fmt.Errorf("gravar binario: %w", err)
	}
	if runtime.GOOS != "windows" {
		if err := os.Chmod(tmp, 0o755); err != nil {
			os.Remove(tmp)
			return fmt.Errorf("chmod binario: %w", err)
		}
	}
	if err := withRetry(5, func() error { return os.Rename(tmp, dest) }); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("instalar binario: %w", err)
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
