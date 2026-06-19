package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"booptube/config"
	"booptube/downloader"
	"booptube/ui"
	"booptube/buildinfo"
)

func main() {
	dirFlag := flag.String("dir", "", "pasta de destino (pula prompt da pasta)")
	versionFlag := flag.Bool("version", false, "mostra versão e sai")
	flag.Parse()

	if *versionFlag {
		fmt.Println("booptube", buildinfo.Info())
		os.Exit(0)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao carregar config: %v\n", err)
		os.Exit(1)
	}
	if *dirFlag != "" {
		cfg.DownloadDir = *dirFlag
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

	if err := ui.Run(ctx, &cfg, dl, *dirFlag != ""); err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		os.Exit(1)
	}
}
