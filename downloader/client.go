package downloader

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"booptube/assets"
	"booptube/config"
	"booptube/video"
)

type Handlers struct {
	OnLine     func(string)
	OnProgress func(float64)
}

var progressRE = regexp.MustCompile(`(\d+(?:\.\d+)?)%`)

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
	if c.cfg.YtdlpPath == "" {
		return fmt.Errorf("caminho do yt-dlp indisponivel")
	}
	if isInstalledBinary(c.cfg.YtdlpPath) {
		return nil
	}
	if len(assets.Ytdlp) == 0 {
		return fmt.Errorf("yt-dlp nao encontrado: reinstale booptube ou compile o build portable com fetch-deps")
	}
	return ensureBinary(ctx, assets.Ytdlp, c.cfg.YtdlpPath)
}

func (c *Client) VerifyTools(ctx context.Context) error {
	if err := c.EnsureYtdlp(ctx); err != nil {
		return err
	}
	if err := c.EnsureFfmpeg(ctx); err != nil {
		return err
	}
	if err := runToolCheck(ctx, c.cfg.YtdlpPath, "--version"); err != nil {
		return fmt.Errorf("yt-dlp: %w", err)
	}
	ffmpegPath := filepath.Join(c.cfg.FfmpegDir, assets.FfmpegName)
	ffprobePath := filepath.Join(c.cfg.FfmpegDir, assets.FfprobeName)
	if err := runToolCheck(ctx, ffmpegPath, "-version"); err != nil {
		return fmt.Errorf("ffmpeg: %w", err)
	}
	if err := runToolCheck(ctx, ffprobePath, "-version"); err != nil {
		return fmt.Errorf("ffprobe: %w", err)
	}
	return nil
}

func runToolCheck(ctx context.Context, path string, args ...string) error {
	cmd := exec.CommandContext(ctx, path, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func (c *Client) EnsureFfmpeg(ctx context.Context) error {
	if c.cfg.FfmpegDir == "" {
		return fmt.Errorf("caminho do ffmpeg indisponivel")
	}
	ffmpegPath := filepath.Join(c.cfg.FfmpegDir, assets.FfmpegName)
	ffprobePath := filepath.Join(c.cfg.FfmpegDir, assets.FfprobeName)
	if isInstalledBinary(ffmpegPath) && isInstalledBinary(ffprobePath) {
		return nil
	}
	if len(assets.Ffmpeg) == 0 || len(assets.Ffprobe) == 0 {
		return fmt.Errorf("ffmpeg nao encontrado: reinstale booptube ou compile o build portable com fetch-deps")
	}
	if err := os.MkdirAll(c.cfg.FfmpegDir, 0o755); err != nil {
		return fmt.Errorf("criar cache ffmpeg: %w", err)
	}
	if err := ensureBinary(ctx, assets.Ffmpeg, ffmpegPath); err != nil {
		return fmt.Errorf("instalar ffmpeg: %w", err)
	}
	if err := ensureBinary(ctx, assets.Ffprobe, ffprobePath); err != nil {
		return fmt.Errorf("instalar ffprobe: %w", err)
	}
	return nil
}

func isInstalledBinary(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	if runtime.GOOS == "windows" {
		return true
	}
	return info.Mode()&0o111 != 0
}

func (c *Client) Download(ctx context.Context, rawURL string, format video.Format, h *Handlers) error {
	if c.cfg.DownloadDir == "" {
		return fmt.Errorf("pasta de destino nao definida")
	}
	_, err := c.DownloadTo(ctx, rawURL, format, c.cfg.DownloadDir, h)
	return err
}

func (c *Client) DownloadTo(ctx context.Context, rawURL string, format video.Format, outDir string, h *Handlers) (string, error) {
	if outDir == "" {
		return "", fmt.Errorf("pasta de destino nao definida")
	}
	if err := ensureWritableDir(outDir); err != nil {
		return "", err
	}
	if err := c.EnsureYtdlp(ctx); err != nil {
		return "", err
	}
	if err := c.EnsureFfmpeg(ctx); err != nil {
		return "", err
	}

	out := filepath.Join(outDir, "%(title)s.%(ext)s")
	args := buildYtdlpArgs(c.cfg.FfmpegDir, out, rawURL, format)
	if args == nil {
		return "", fmt.Errorf("formato nao suportado: %s", format)
	}
	cmd := exec.CommandContext(ctx, c.cfg.YtdlpPath, args...)
	if h == nil {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("download falhou: %w", err)
		}
		return findOutputFile(outDir)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("stderr pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("iniciar yt-dlp: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		scanStream(stdout, h)
	}()
	go func() {
		defer wg.Done()
		scanStream(stderr, h)
	}()
	wg.Wait()

	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("download falhou: %w", err)
	}
	return findOutputFile(outDir)
}

func buildYtdlpArgs(ffmpegDir, out, rawURL string, format video.Format) []string {
	args := []string{
		"--no-playlist",
		"--ffmpeg-location", ffmpegDir,
		"-o", out,
		rawURL,
	}
	switch format {
	case video.FormatMP4:
		return append([]string{
			"-f", "bv*+ba/b",
			"--merge-output-format", "mp4",
		}, args...)
	case video.FormatMP3:
		return append([]string{
			"-x",
			"--audio-format", "mp3",
		}, args...)
	default:
		return nil
	}
}

func findOutputFile(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("listar pasta destino: %w", err)
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".part") || strings.HasSuffix(name, ".tmp") || strings.HasPrefix(name, ".") {
			continue
		}
		files = append(files, filepath.Join(dir, name))
	}
	if len(files) == 0 {
		return "", fmt.Errorf("nenhum arquivo gerado em %s", dir)
	}
	if len(files) == 1 {
		return files[0], nil
	}
	var best string
	var bestMod time.Time
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			continue
		}
		if info.ModTime().After(bestMod) {
			bestMod = info.ModTime()
			best = f
		}
	}
	if best == "" {
		return "", fmt.Errorf("nenhum arquivo valido em %s", dir)
	}
	return best, nil
}

func scanStream(r io.Reader, h *Handlers) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		if h.OnLine != nil {
			h.OnLine(line)
		}
		if h.OnProgress != nil {
			if m := progressRE.FindStringSubmatch(line); len(m) > 1 {
				var pct float64
				if _, err := fmt.Sscanf(m[1], "%f", &pct); err == nil {
					h.OnProgress(pct)
				}
			}
		}
	}
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
