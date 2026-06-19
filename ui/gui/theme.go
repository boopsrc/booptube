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
	neonBlue       = color.NRGBA{R: 37, G: 99, B: 235, A: 255}
	neonBlueHover  = color.NRGBA{R: 29, G: 78, B: 216, A: 255}
	neonBluePress  = color.NRGBA{R: 30, G: 64, B: 175, A: 255}
	neonCardBorder = color.NRGBA{R: 0, G: 240, B: 255, A: 80}
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
		return color.NRGBA{R: 30, G: 58, B: 95, A: 255}
	case theme.ColorNameDisabled:
		return neonMuted
	case theme.ColorNameError:
		return neonError
	case theme.ColorNameForeground:
		return neonText
	case theme.ColorNameForegroundOnPrimary:
		return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 24, G: 24, B: 40, A: 255}
	case theme.ColorNameInputBorder:
		return neonCardBorder
	case theme.ColorNamePlaceHolder:
		return neonMuted
	case theme.ColorNamePrimary:
		return neonBlue
	case theme.ColorNameHover:
		return neonBlueHover
	case theme.ColorNamePressed:
		return neonBluePress
	case theme.ColorNameScrollBar:
		return neonSecondary
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 60}
	case theme.ColorNameSelection:
		return color.NRGBA{R: 37, G: 99, B: 235, A: 80}
	case theme.ColorNameSeparator:
		return neonCardBorder
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
		return 16
	case theme.SizeNameInnerPadding:
		return 10
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 28
	default:
		return t.base.Size(name)
	}
}
