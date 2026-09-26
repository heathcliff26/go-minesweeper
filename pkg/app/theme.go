package app

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const colorNameTileDefault = "tileDefault"

var (
	lightGray                    = color.RGBA{211, 211, 211, alpha}
	tileDefaultColorVariantLight = color.RGBA{180, 180, 180, alpha}
	tileDefaultColorVariantDark  = color.RGBA{100, 100, 100, alpha}

	borderShadowVariantLight = color.Black
	borderShadowVariantDark  = color.RGBA{120, 120, 120, alpha}
)

var (
	// Ensure there are compile errors if the theme interface is not implemented
	_ fyne.Theme = mainTheme{}
	_ fyne.Theme = borderTheme{}
)

type mainTheme struct{}

func (mainTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if variant == theme.VariantLight && name == theme.ColorNameBackground {
		return lightGray
	}
	if name == colorNameTileDefault {
		switch variant {
		case theme.VariantLight:
			return tileDefaultColorVariantLight
		case theme.VariantDark:
			return tileDefaultColorVariantDark
		}
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (mainTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (mainTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (mainTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

type borderTheme struct {
	mainTheme
}

func (borderTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameShadow {
		switch variant {
		case theme.VariantLight:
			return borderShadowVariantLight
		case theme.VariantDark:
			return borderShadowVariantDark
		}
	}
	return mainTheme{}.Color(name, variant)
}
