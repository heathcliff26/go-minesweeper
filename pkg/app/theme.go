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
	tileDefaultColorVariantDark  = color.Gray16{32767}
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
		if variant == theme.VariantLight {
			return tileDefaultColorVariantLight
		} else if variant == theme.VariantDark {
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
		if variant == theme.VariantLight {
			return color.Black
		} else if variant == theme.VariantDark {
			return color.White
		}
	}
	return mainTheme{}.Color(name, variant)
}
