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
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const maxLogLines = 20

func Run(ctx context.Context, cfg *config.Config, dl *downloader.Client) error {
	a := app.NewWithID("dev.booptube.gui")
	a.Settings().SetTheme(newNeonTheme())

	w := a.NewWindow("booptube")
	w.Resize(fyne.NewSize(720, 640))
	w.SetFixedSize(false)

	dirEntry := widget.NewEntry()
	dirEntry.SetPlaceHolder("C:\\Users\\voce\\Downloads")
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
	statusLabel.Importance = widget.MediumImportance

	progressBar := widget.NewProgressBar()
	progressBar.Min = 0
	progressBar.Max = 100

	logEntry := widget.NewMultiLineEntry()
	logEntry.SetPlaceHolder("Saída do yt-dlp aparecerá aqui...")
	logEntry.Disable()
	logEntry.Wrapping = fyne.TextWrapWord
	logEntry.SetMinRowsVisible(5)

	downloadBtn := widget.NewButton("Baixar", nil)
	downloadBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButton("Cancelar", nil)
	cancelBtn.Importance = widget.MediumImportance
	cancelBtn.Hide()

	chooseDirBtn := widget.NewButton("Escolher...", nil)

	var (
		mu             sync.Mutex
		logLines       []string
		downloadCancel context.CancelFunc
		downloading    bool
	)

	setStatus := func(msg string, imp widget.Importance) {
		fyne.Do(func() {
			statusLabel.SetText(msg)
			statusLabel.Importance = imp
		})
	}

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
				chooseDirBtn.Enable()
				downloadBtn.Enable()
			} else {
				dirEntry.Disable()
				urlEntry.Disable()
				formatRadio.Disable()
				chooseDirBtn.Disable()
				downloadBtn.Disable()
			}
		})
	}

	showDownloadUI := func(active bool) {
		downloading = active
		fyne.Do(func() {
			if active {
				progressBar.SetValue(0)
				cancelBtn.Show()
				cancelBtn.Enable()
			} else {
				progressBar.SetValue(0)
				cancelBtn.Hide()
			}
		})
	}

	finishDownload := func(err error, savedDir string) {
		fyne.Do(func() {
			showDownloadUI(false)
			setInputsEnabled(true)
			if err != nil {
				if err == context.Canceled {
					statusLabel.SetText("Download cancelado.")
					statusLabel.Importance = widget.MediumImportance
				} else {
					statusLabel.SetText(fmt.Sprintf("Erro: %v", err))
					statusLabel.Importance = widget.DangerImportance
				}
				return
			}
			cfg.DownloadDir = savedDir
			if saveErr := config.Save(*cfg); saveErr != nil {
				statusLabel.SetText(fmt.Sprintf("Concluído (aviso: config não salva: %v)", saveErr))
				statusLabel.Importance = widget.WarningImportance
			} else {
				statusLabel.SetText("Concluído.")
				statusLabel.Importance = widget.SuccessImportance
			}
			urlEntry.SetText("")
		})
	}

	startDownload := func() {
		dir := strings.TrimSpace(dirEntry.Text)
		if dir == "" {
			setStatus("Informe a pasta de destino.", widget.WarningImportance)
			return
		}
		if err := ui.EnsureDir(dir); err != nil {
			setStatus(fmt.Sprintf("Erro: %v", err), widget.DangerImportance)
			return
		}

		rawURL := strings.TrimSpace(urlEntry.Text)
		parsed, err := video.ParseURL(rawURL)
		if err != nil {
			setStatus(fmt.Sprintf("Erro: %v", err), widget.DangerImportance)
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
		setStatus(fmt.Sprintf("Baixando %s como %s...", parsed, format), widget.MediumImportance)

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

	chooseDirBtn.OnTapped = func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				setStatus(fmt.Sprintf("Erro: %v", err), widget.DangerImportance)
				return
			}
			if uri == nil {
				return
			}
			dir := uri.Path()
			if err := ui.EnsureDir(dir); err != nil {
				setStatus(fmt.Sprintf("Erro: %v", err), widget.DangerImportance)
				return
			}
			dirEntry.SetText(dir)
		}, w)
	}

	w.SetCloseIntercept(func() {
		if downloading && downloadCancel != nil {
			downloadCancel()
		}
		w.Close()
	})

	subtitle := canvas.NewText("YouTube → MP4 / MP3", neonMuted)
	subtitle.TextSize = 14

	header := container.NewVBox(
		neonGlowTitle("booptube"),
		subtitle,
	)

	dirRow := container.NewBorder(nil, nil, nil, chooseDirBtn, dirEntry)
	formContent := container.NewVBox(
		fieldLabel("Pasta de destino"),
		dirRow,
		fieldLabel("URL do YouTube"),
		urlEntry,
		fieldLabel("Formato"),
		formatRadio,
		container.NewBorder(nil, nil, cancelBtn, downloadBtn, layout.NewSpacer()),
	)

	progressContent := container.NewVBox(
		statusLabel,
		progressBar,
	)

	logScroll := container.NewScroll(logEntry)
	logScroll.SetMinSize(fyne.NewSize(0, 120))

	content := container.NewVBox(
		header,
		sectionCard("Download", formContent),
		sectionCard("Progresso", progressContent),
		sectionCard("Log", logScroll),
	)

	bg := canvas.NewRectangle(neonBackground)
	w.SetContent(container.NewStack(bg, container.NewPadded(content)))

	w.Show()
	a.Run()
	return ctx.Err()
}

func fieldLabel(text string) fyne.CanvasObject {
	l := widget.NewLabel(text)
	l.TextStyle = fyne.TextStyle{Bold: true}
	return l
}

func sectionCard(title string, body fyne.CanvasObject) fyne.CanvasObject {
	titleText := canvas.NewText(title, neonPrimary)
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.TextSize = 13

	inner := container.NewVBox(titleText, body)
	padded := container.NewPadded(inner)

	bg := canvas.NewRectangle(neonSurface)
	border := canvas.NewRectangle(neonCardBorder)
	border.StrokeWidth = 1
	border.StrokeColor = neonCardBorder
	border.FillColor = color.Transparent

	return container.NewStack(border, bg, padded)
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
