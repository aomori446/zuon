package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/aomori446/zuon/front/i18n"
)

func Start() {
	a := app.NewWithID("com.aomori446.zuon")
	i18n.Init()
	
	w := a.NewWindow(i18n.T("app_title"))
	
	refreshWindow(w)
	
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}

func refreshWindow(w fyne.Window) {
	w.SetTitle(i18n.T("app_title"))
	
	w.SetMainMenu(nil)
	
	langMap := map[string]string{
		"中文":    "zh",
		"English": "en",
		"日本語":  "ja",
	}
	
	codeMap := map[string]string{
		"zh": "中文",
		"en": "English",
		"ja": "日本語",
	}
	
	langSelect := widget.NewSelect([]string{"中文", "English", "日本語"}, func(s string) {
		if code, ok := langMap[s]; ok {
			
			if code != i18n.GetLang() {
				i18n.SetLang(code)
				refreshWindow(w)
			}
		}
	})
	
	if label, ok := codeMap[i18n.GetLang()]; ok {
		langSelect.SetSelected(label)
	} else {
		langSelect.SetSelected("中文")
	}
	
	topBar := container.NewHBox(
		layout.NewSpacer(),
		widget.NewLabel(i18n.T("menu_language")+":"),
		langSelect,
	)
	
	tabs := container.NewAppTabs()
	tabs.SetTabLocation(container.TabLocationLeading)
	
	tabs.Append(NewEmbedTab(w))
	tabs.Append(NewExtractTab(w))
	
	content := container.NewBorder(container.NewPadded(topBar), nil, nil, nil, tabs)
	
	w.SetContent(content)
}
