package gui

import (
	"context"
	"fmt"
	"image/color"
	"strings"
	"sync"

	"booptube/config"
	"booptube/downloader"
	"booptube/ui"
	"booptube/video"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const maxLogLines = 20

func Run(ctx context.Context, cfg *config.Config, dl *downloader.Client) error {
	a := app.NewWithID("dev.booptube.gui")
	a.Settings().SetTheme(newNeonTheme())

	w := a.NewWindow("booptube")
	w.Resize(fyne.NewSize(680, 560))
	w.SetFixedSize(false)

	dirEntry := widget.NewEntry()
	dirEntry.SetPlaceHolder("Pasta de destino")
	if cfg.DownloadDir != "" {
		dirEntry.SetText(cfg.DownloadDir)
	}

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("https://youtube.com/watch?v=...")

	formatRadio := widget.NewRadioGroup([]string{"MP4 (vídeo)", "MP3 (áudio)"}, nil)
	formatRadio.SetSelected("MP4 (vídeo)")
	formatRadio.Horizontal = true

	statusLabel := widget.NewLabel("Pronto.")
	statusLabel.Wrapping = fyne.TextWrapWord

	progressBar := widget.NewProgressBar()
	progressBar.Min = 0
	progressBar.Max = 100
	progressBar.Hide()

	logEntry := widget.NewMultiLineEntry()
	logEntry.SetPlaceHolder("Log do download...")
	logEntry.Disable()
	logEntry.Wrapping = fyne.TextWrapWord
	logEntry.SetMinRowsVisible(6)

	downloadBtn := widget.NewButton("Baixar", nil)
	cancelBtn := widget.NewButton("Cancelar", nil)
	cancelBtn.Hide()

	var (
		mu             sync.Mutex
		logLines       []string
		downloadCancel context.CancelFunc
		downloading    bool
	)

	appendLog := func(line string) {
		mu.Lock()
		logLines = append(logLines, line)
		if len(logLines) > maxLogLines {
			logLines = logLines[len(logLines)-maxLogLines:]
		}
		text := strings.Join(logLines, "\n")
		mu.Unlock()
		fyne.Do(func() {
			logEntry.SetText(text)
		})
	}

	setStatus := func(msg string) {
		fyne.Do(func() {
			statusLabel.SetText(msg)
		})
	}

	setProgress := func(pct float64) {
		fyne.Do(func() {
			progressBar.SetValue(pct)
		})
	}

	setInputsEnabled := func(enabled bool) {
		fyne.Do(func() {
			if enabled {
				dirEntry.Enable()
				urlEntry.Enable()
				formatRadio.Enable()
				downloadBtn.Enable()
			} else {
				dirEntry.Disable()
				urlEntry.Disable()
				formatRadio.Disable()
				downloadBtn.Disable()
			}
		})
	}

	showDownloadUI := func(active bool) {
		downloading = active
		if active {
			progressBar.Show()
			progressBar.SetValue(0)
			cancelBtn.Show()
		} else {
			progressBar.Hide()
			cancelBtn.Hide()
		}
	}

	finishDownload := func(err error, savedDir string) {
		fyne.Do(func() {
			showDownloadUI(false)
			setInputsEnabled(true)
			if err != nil {
				if err == context.Canceled {
					statusLabel.SetText("Download cancelado.")
				} else {
					statusLabel.SetText(fmt.Sprintf("Erro: %v", err))
				}
				return
			}
			cfg.DownloadDir = savedDir
			if saveErr := config.Save(*cfg); saveErr != nil {
				statusLabel.SetText(fmt.Sprintf("Concluído (aviso: config não salva: %v)", saveErr))
			} else {
				statusLabel.SetText("Concluído.")
			}
			urlEntry.SetText("")
		})
	}

	startDownload := func() {
		dir := strings.TrimSpace(dirEntry.Text)
		if dir == "" {
			setStatus("Informe a pasta de destino.")
			return
		}
		if err := ui.EnsureDir(dir); err != nil {
			setStatus(fmt.Sprintf("Erro: %v", err))
			return
		}

		rawURL := strings.TrimSpace(urlEntry.Text)
		parsed, err := video.ParseURL(rawURL)
		if err != nil {
			setStatus(fmt.Sprintf("Erro: %v", err))
			return
		}

		format := video.FormatMP4
		if formatRadio.Selected == "MP3 (áudio)" {
			format = video.FormatMP3
		}

		mu.Lock()
		logLines = nil
		mu.Unlock()
		logEntry.SetText("")

		dlCfg := dl.Config()
		dlCfg.DownloadDir = dir
		*dl = *downloader.New(dlCfg)

		dlCtx, cancel := context.WithCancel(ctx)
		downloadCancel = cancel

		setInputsEnabled(false)
		showDownloadUI(true)
		setStatus(fmt.Sprintf("Baixando %s como %s...", parsed, format))

		handlers := &downloader.Handlers{
			OnLine: appendLog,
			OnProgress: func(pct float64) {
				setProgress(pct)
			},
		}

		go func() {
			defer func() {
				if r := recover(); r != nil {
					finishDownload(fmt.Errorf("panic: %v", r), dir)
				}
			}()
			err := dl.Download(dlCtx, parsed, format, handlers)
			finishDownload(err, dir)
		}()
	}

	downloadBtn.OnTapped = startDownload

	cancelBtn.OnTapped = func() {
		if downloadCancel != nil {
			downloadCancel()
		}
	}

	chooseDirBtn := widget.NewButton("Escolher...", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				setStatus(fmt.Sprintf("Erro: %v", err))
				return
			}
			if uri == nil {
				return
			}
			dir := uri.Path()
			if err := ui.EnsureDir(dir); err != nil {
				setStatus(fmt.Sprintf("Erro: %v", err))
				return
			}
			dirEntry.SetText(dir)
		}, w)
	})

	w.SetCloseIntercept(func() {
		if downloading && downloadCancel != nil {
			downloadCancel()
		}
		w.Close()
	})

	header := neonGlowTitle("booptube")

	dirRow := container.NewBorder(nil, nil, nil, chooseDirBtn, dirEntry)
	dirLabel := widget.NewLabel("Pasta de destino")
	urlLabel := widget.NewLabel("URL do YouTube")
	formatLabel := widget.NewLabel("Formato")

	btnRow := container.NewHBox(downloadBtn, cancelBtn)

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		dirLabel,
		dirRow,
		urlLabel,
		urlEntry,
		formatLabel,
		formatRadio,
		btnRow,
		statusLabel,
		progressBar,
		widget.NewLabel("Log"),
		logEntry,
	)

	scroll := container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(640, 520))

	bg := canvas.NewRectangle(neonBackground)
	w.SetContent(container.NewStack(bg, container.NewPadded(scroll)))

	w.Show()
	a.Run()
	return ctx.Err()
}

func neonGlowTitle(text string) fyne.CanvasObject {
	layers := []struct {
		offset fyne.Position
		alpha  uint8
	}{
		{fyne.NewPos(-2, 0), 30},
		{fyne.NewPos(2, 0), 30},
		{fyne.NewPos(0, -2), 30},
		{fyne.NewPos(0, 2), 30},
		{fyne.NewPos(-1, -1), 50},
		{fyne.NewPos(1, 1), 50},
		{fyne.NewPos(0, 0), 255},
	}

	var objs []fyne.CanvasObject
	for _, layer := range layers {
		c := neonPrimary
		if layer.alpha < 255 {
			c = color.NRGBA{R: 0, G: 240, B: 255, A: layer.alpha}
		}
		t := canvas.NewText(text, c)
		t.TextStyle = fyne.TextStyle{Bold: true}
		t.TextSize = 32
		t.Move(layer.offset)
		objs = append(objs, t)
	}
	return container.NewStack(objs...)
}
