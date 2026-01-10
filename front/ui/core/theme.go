package core

import (
	"embed"
	"image/color"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/aomori446/zuon/front/i18n"
)

//go:embed assets/fonts/*.ttf
var fontAssets embed.FS

const (
	ThemeModeSystem int = 0
	ThemeModeDark   int = 1
	ThemeModeLight  int = 2
	ThemeModeOcean  int = 3
	ThemeModeForest int = 4
)

type myTheme struct {
	mode int
}

func NewMyTheme(mode int) fyne.Theme {
	return &myTheme{mode: mode}
}

var _ fyne.Theme = (*myTheme)(nil)

func (m *myTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch m.mode {
	case ThemeModeOcean:
		switch n {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 15, G: 23, B: 42, A: 255} // Slate 900
		case theme.ColorNameInputBackground, theme.ColorNameOverlayBackground, theme.ColorNameMenuBackground:
			return color.NRGBA{R: 30, G: 41, B: 59, A: 255} // Slate 800
		case theme.ColorNameForeground, theme.ColorNamePlaceHolder:
			return color.NRGBA{R: 226, G: 232, B: 240, A: 255} // Slate 200
		case theme.ColorNamePrimary, theme.ColorNameHyperlink:
			return color.NRGBA{R: 56, G: 189, B: 248, A: 255} // Sky 400
		case theme.ColorNameButton:
			return color.NRGBA{R: 2, G: 132, B: 199, A: 255} // Sky 600
		case theme.ColorNameFocus:
			return color.NRGBA{R: 56, G: 189, B: 248, A: 120} // Sky 400 Transparent
		case theme.ColorNameSelection:
			return color.NRGBA{R: 51, G: 65, B: 85, A: 255} // Slate 700
		case theme.ColorNameScrollBar:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 30}
		case theme.ColorNameShadow:
			return color.NRGBA{R: 0, G: 0, B: 0, A: 150}
		default:
			// Fallback to Dark theme
			return theme.DefaultTheme().Color(n, theme.VariantDark)
		}
	case ThemeModeForest:
		switch n {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 20, G: 24, B: 20, A: 255} // Deep Dark Green
		case theme.ColorNameInputBackground, theme.ColorNameOverlayBackground, theme.ColorNameMenuBackground:
			return color.NRGBA{R: 40, G: 48, B: 40, A: 255} // Dark Olive
		case theme.ColorNameForeground, theme.ColorNamePlaceHolder:
			return color.NRGBA{R: 220, G: 225, B: 215, A: 255} // Soft Beige
		case theme.ColorNamePrimary, theme.ColorNameHyperlink:
			return color.NRGBA{R: 150, G: 200, B: 130, A: 255} // Matcha Green
		case theme.ColorNameButton:
			return color.NRGBA{R: 60, G: 90, B: 60, A: 255} // Forest Green
		case theme.ColorNameFocus:
			return color.NRGBA{R: 150, G: 200, B: 130, A: 120} // Matcha Transparent
		case theme.ColorNameSelection:
			return color.NRGBA{R: 70, G: 90, B: 70, A: 255} // Lighter Olive
		case theme.ColorNameScrollBar:
			return color.NRGBA{R: 255, G: 255, B: 255, A: 30}
		case theme.ColorNameShadow:
			return color.NRGBA{R: 0, G: 0, B: 0, A: 150}
		default:
			// Fallback to Dark theme
			return theme.DefaultTheme().Color(n, theme.VariantDark)
		}
	}

	targetVariant := v
	if m.mode == ThemeModeDark {
		targetVariant = theme.VariantDark
	} else if m.mode == ThemeModeLight {
		targetVariant = theme.VariantLight
	}
	
	return theme.DefaultTheme().Color(n, targetVariant)
}

func (m *myTheme) Font(s fyne.TextStyle) fyne.Resource {
	lang := i18n.GetLang()
	
	if lang == "my" {
		if s.Bold {
			if data, err := fontAssets.ReadFile("assets/fonts/NotoSansMyanmar-Bold.ttf"); err == nil {
				return fyne.NewStaticResource("NotoSansMyanmar-Bold.ttf", data)
			}
		}
		if data, err := fontAssets.ReadFile("assets/fonts/NotoSansMyanmar-Regular.ttf"); err == nil {
			return fyne.NewStaticResource("NotoSansMyanmar-Regular.ttf", data)
		}
	}
	
	// Fallback to default theme font for other languages
	return theme.DefaultTheme().Font(s)
}

func (m *myTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (m *myTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}

func ApplyTheme(a fyne.App) {
	mode := a.Preferences().Int("theme_mode")
	a.Settings().SetTheme(NewMyTheme(mode))
}