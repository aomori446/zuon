package ui

import (
	"embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/aomori446/zuon/front/i18n"
)

//go:embed assets/fonts/*.ttf
var fontAssets embed.FS

type myTheme struct {
}

var _ fyne.Theme = (*myTheme)(nil)

func (m *myTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
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

	// Fallback to default theme font
	return theme.DefaultTheme().Font(s)
}

func (m *myTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (m *myTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
