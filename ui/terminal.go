package ui

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"booptube/config"
	"booptube/downloader"
	"booptube/video"
)

func Run(ctx context.Context, cfg *config.Config, dl *downloader.Client, skipDirPrompt bool) error {
	sc := bufio.NewScanner(os.Stdin)
	fmt.Fprintln(os.Stderr, "booptube — digite q ou sair para encerrar")

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		var dir string
		if skipDirPrompt && cfg.DownloadDir != "" {
			dir = cfg.DownloadDir
			if err := ensureDir(dir); err != nil {
				fmt.Fprintf(os.Stderr, "erro: %v\n", err)
				skipDirPrompt = false
				continue
			}
		} else {
			var err error
			dir, err = promptDir(sc, cfg.DownloadDir)
			if err != nil {
				return err
			}
			if dir == "" {
				return nil
			}
		}
		cfg.DownloadDir = dir

		rawURL, err := promptLine(sc, "URL do YouTube: ")
		if err != nil {
			return err
		}
		if isQuit(rawURL) {
			return nil
		}
		parsed, err := video.ParseURL(rawURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erro: %v\n", err)
			continue
		}

		fmtRaw, err := promptLine(sc, "Formato [1=mp4, 2=mp3] (Enter=mp4): ")
		if err != nil {
			return err
		}
		if isQuit(fmtRaw) {
			return nil
		}
		if strings.TrimSpace(fmtRaw) == "" {
			fmtRaw = "1"
		}
		format, err := video.FormatFromString(fmtRaw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erro: %v\n", err)
			continue
		}

		dlCfg := dl.Config()
		dlCfg.DownloadDir = cfg.DownloadDir
		*dl = *downloader.New(dlCfg)

		fmt.Fprintf(os.Stderr, "baixando %s como %s...\n", parsed, format)
		if err := dl.Download(ctx, parsed, format, nil); err != nil {
			fmt.Fprintf(os.Stderr, "erro: %v\n", err)
			continue
		}
		fmt.Fprintln(os.Stderr, "concluido.")

		if err := config.Save(*cfg); err != nil {
			fmt.Fprintf(os.Stderr, "aviso: nao foi possivel salvar config: %v\n", err)
		}
	}
}

func promptDir(sc *bufio.Scanner, current string) (string, error) {
	hint := ""
	if current != "" {
		hint = fmt.Sprintf(" (Enter=%s)", current)
	}
	line, err := promptLine(sc, "Pasta de destino"+hint+": ")
	if err != nil {
		return "", err
	}
	if isQuit(line) {
		return "", nil
	}
	if strings.TrimSpace(line) == "" {
		if current == "" {
			fmt.Fprintln(os.Stderr, "informe uma pasta de destino.")
			return promptDir(sc, current)
		}
		line = current
	}
	dir := filepath.Clean(strings.TrimSpace(line))
	if err := ensureDir(dir); err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		return promptDir(sc, current)
	}
	return dir, nil
}

func promptLine(sc *bufio.Scanner, label string) (string, error) {
	fmt.Fprint(os.Stderr, label)
	if !sc.Scan() {
		if err := sc.Err(); err != nil {
			return "", err
		}
		return "", nil
	}
	return strings.TrimSpace(sc.Text()), nil
}

func ensureDir(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("criar pasta: %w", err)
	}
	test := filepath.Join(dir, ".booptube-write-test")
	if err := os.WriteFile(test, []byte("ok"), 0o644); err != nil {
		return fmt.Errorf("pasta nao gravavel: %w", err)
	}
	return os.Remove(test)
}

func isQuit(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "q", "quit", "sair", "exit":
		return true
	default:
		return false
	}
}
