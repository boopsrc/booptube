package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"booptube/buildinfo"
	"booptube/config"
	"booptube/downloader"
	"booptube/ui/web"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg, err := config.LoadWeb()
	if err != nil {
		slog.Error("config_error", "error", err.Error())
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := os.MkdirAll(cfg.DownloadDir, 0o755); err != nil {
		slog.Error("mkdir_error", "path", cfg.DownloadDir, "error", err.Error())
		os.Exit(1)
	}

	dl := downloader.New(cfg.ToConfig())
	initCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	if err := dl.VerifyTools(initCtx); err != nil {
		cancel()
		slog.Error("tools_verify_error", "error", err.Error())
		os.Exit(1)
	}
	cancel()

	srv := web.New(cfg, dl)
	httpSrv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	slog.Info("server_start",
		"addr", cfg.Addr,
		"version", buildinfo.Version,
		"max_concurrent", cfg.MaxConcurrent,
		"file_ttl", cfg.FileTTL.String(),
		"download_dir", cfg.DownloadDir,
	)

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server_error", "error", err.Error())
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("server_shutdown")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	srv.Shutdown()
	httpSrv.Shutdown(shutdownCtx)
}
