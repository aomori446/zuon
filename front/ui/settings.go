package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/aomori446/zuon/front/i18n"
)

type SettingsTab struct {
	Content fyne.CanvasObject
}

func NewSettingsTab(w fyne.Window, a fyne.App, onRefresh func()) *SettingsTab {
	// --- Language Section ---
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
				applyTheme(a) // Re-apply theme to ensure font update
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

	// --- Theme Section ---
	themeOptions := []string{i18n.T("settings_theme_system"), i18n.T("settings_theme_dark"), i18n.T("settings_theme_light")}
	themeSelect := widget.NewSelect(themeOptions, func(s string) {
		var mode int
		switch s {
		case i18n.T("settings_theme_dark"):
			mode = ThemeModeDark
		case i18n.T("settings_theme_light"):
			mode = ThemeModeLight
		default:
			mode = ThemeModeSystem
		}
		a.Preferences().SetInt("theme_mode", mode)
		applyTheme(a)
	})
	
	// Set current selection
	currentMode := a.Preferences().Int("theme_mode")
	switch currentMode {
	case ThemeModeDark:
		themeSelect.SetSelected(i18n.T("settings_theme_dark"))
	case ThemeModeLight:
		themeSelect.SetSelected(i18n.T("settings_theme_light"))
	default:
		themeSelect.SetSelected(i18n.T("settings_theme_system"))
	}

	themeCard := widget.NewCard(i18n.T("settings_theme"), "", container.NewVBox(themeSelect))

	// --- Account Section ---
	logoutBtn := widget.NewButtonWithIcon(i18n.T("btn_logout"), theme.LogoutIcon(), func() {
		a.Preferences().SetString("auth_token", "")
		w.Close()
		ShowLoginWindow(a, func(token string) {
			a.Preferences().SetString("auth_token", token)
			showMainWindow(a)
		})
	})
	logoutBtn.Importance = widget.HighImportance

	accountCard := widget.NewCard(i18n.T("settings_account"), "", container.NewVBox(logoutBtn))

	// Layout
	content := container.NewVBox(
		langCard,
		themeCard,
		accountCard,
	)

	return &SettingsTab{
		Content: container.NewPadded(container.NewVScroll(content)),
	}
}
