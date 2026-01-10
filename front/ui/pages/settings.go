package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/aomori446/zuon/front/i18n"
	"github.com/aomori446/zuon/front/ui/core"
)

type SettingsTab struct {
	Content fyne.CanvasObject
}

func NewSettingsTab(a fyne.App, onRefresh func(), onLogout func()) *SettingsTab {
	langMap := map[string]string{
		"Chinese":  "zh",
		"English":  "en",
		"Japanese": "ja",
		"Myanmar":  "my",
	}
	
	codeMap := map[string]string{
		"zh": "Chinese",
		"en": "English",
		"ja": "Japanese",
		"my": "Myanmar",
	}
	
	langSelect := widget.NewSelect([]string{"Chinese", "English", "Japanese", "Myanmar"}, func(s string) {
		if code, ok := langMap[s]; ok {
			if code != i18n.GetLang() {
				i18n.SetLang(code)
				core.ApplyTheme(a) // Re-apply theme to ensure font update
				onRefresh()
			}
		}
	})
	if label, ok := codeMap[i18n.GetLang()]; ok {
		langSelect.SetSelected(label)
	} else {
		langSelect.SetSelected("English")
	}
	
	langCard := widget.NewCard(i18n.T("settings_language"), "", container.NewVBox(langSelect))
	
	themeOptions := []string{
		i18n.T("settings_theme_system"),
		i18n.T("settings_theme_dark"),
		i18n.T("settings_theme_light"),
		i18n.T("settings_theme_ocean"),
		i18n.T("settings_theme_forest"),
	}
	themeSelect := widget.NewSelect(themeOptions, func(s string) {
		var mode int
		switch s {
		case i18n.T("settings_theme_dark"):
			mode = core.ThemeModeDark
		case i18n.T("settings_theme_light"):
			mode = core.ThemeModeLight
		case i18n.T("settings_theme_ocean"):
			mode = core.ThemeModeOcean
		case i18n.T("settings_theme_forest"):
			mode = core.ThemeModeForest
		default:
			mode = core.ThemeModeSystem
		}
		a.Preferences().SetInt("theme_mode", mode)
		core.ApplyTheme(a)
	})
	
	currentMode := a.Preferences().Int("theme_mode")
	switch currentMode {
	case core.ThemeModeDark:
		themeSelect.SetSelected(i18n.T("settings_theme_dark"))
	case core.ThemeModeLight:
		themeSelect.SetSelected(i18n.T("settings_theme_light"))
	case core.ThemeModeOcean:
		themeSelect.SetSelected(i18n.T("settings_theme_ocean"))
	case core.ThemeModeForest:
		themeSelect.SetSelected(i18n.T("settings_theme_forest"))
	default:
		themeSelect.SetSelected(i18n.T("settings_theme_system"))
	}
	
	themeCard := widget.NewCard(i18n.T("settings_theme"), "", container.NewVBox(themeSelect))
	
	logoutBtn := widget.NewButtonWithIcon(i18n.T("btn_logout"), theme.LogoutIcon(), func() {
		if onLogout != nil {
			onLogout()
		}
	})
	logoutBtn.Importance = widget.HighImportance
	
	accountCard := widget.NewCard(i18n.T("settings_account"), "", container.NewVBox(logoutBtn))
	
	content := container.NewVBox(
		langCard,
		themeCard,
		accountCard,
	)
	
	return &SettingsTab{
		Content: container.NewPadded(container.NewVScroll(content)),
	}
}
