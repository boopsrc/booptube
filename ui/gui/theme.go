package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	neonBackground = color.NRGBA{R: 10, G: 10, B: 18, A: 255}
	neonSurface    = color.NRGBA{R: 18, G: 18, B: 31, A: 255}
	neonPrimary    = color.NRGBA{R: 0, G: 240, B: 255, A: 255}
	neonSecondary  = color.NRGBA{R: 191, G: 0, B: 255, A: 255}
	neonError      = color.NRGBA{R: 255, G: 51, B: 102, A: 255}
	neonText       = color.NRGBA{R: 224, G: 232, B: 255, A: 255}
	neonMuted      = color.NRGBA{R: 107, G: 114, B: 128, A: 255}
)

type neonTheme struct {
	base fyne.Theme
}

func newNeonTheme() fyne.Theme {
	return &neonTheme{base: theme.DefaultTheme()}
}

func (t *neonTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return neonBackground
	case theme.ColorNameButton:
		return neonSurface
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 30, G: 30, B: 45, A: 255}
	case theme.ColorNameDisabled:
		return neonMuted
	case theme.ColorNameError:
		return neonError
	case theme.ColorNameForeground:
		return neonText
	case theme.ColorNameInputBackground:
		return neonSurface
	case theme.ColorNameInputBorder:
		return neonPrimary
	case theme.ColorNamePlaceHolder:
		return neonMuted
	case theme.ColorNamePrimary:
		return neonPrimary
	case theme.ColorNameHover:
		return color.NRGBA{R: 0, G: 200, B: 220, A: 255}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 0, G: 160, B: 180, A: 255}
	case theme.ColorNameScrollBar:
		return neonSecondary
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 240, B: 255, A: 40}
	case theme.ColorNameSelection:
		return color.NRGBA{R: 191, G: 0, B: 255, A: 80}
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 0, G: 240, B: 255, A: 60}
	default:
		return t.base.Color(name, variant)
	}
}

func (t *neonTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.base.Font(style)
}

func (t *neonTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.base.Icon(name)
}

func (t *neonTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 12
	case theme.SizeNameInnerPadding:
		return 8
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 28
	default:
		return t.base.Size(name)
	}
}
