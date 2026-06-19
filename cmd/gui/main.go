package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"booptube/config"
	"booptube/downloader"
	"booptube/ui/gui"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao carregar config: %v\n", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dl := downloader.New(cfg)
	initCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	if err := dl.EnsureYtdlp(initCtx); err != nil {
		cancel()
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		os.Exit(1)
	}
	if err := dl.EnsureFfmpeg(initCtx); err != nil {
		cancel()
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		os.Exit(1)
	}
	cancel()

	if err := gui.Run(ctx, &cfg, dl); err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		os.Exit(1)
	}
}
