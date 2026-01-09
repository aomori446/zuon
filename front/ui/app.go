package ui

import (
	"embed"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/aomori446/zuon/front/i18n"
)

//go:embed assets/*.png
var iconData embed.FS

func Start() {
	a := app.NewWithID("com.aomori446.zuon")
	i18n.Init()
	
	if i18n.GetLang() == "my" {
		a.Settings().SetTheme(&myTheme{})
	} else {
		a.Settings().SetTheme(theme.DefaultTheme())
	}
	
	w := a.NewWindow(i18n.T("app_title"))
	
	refreshWindow(w)
	
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(900, 600))
	w.ShowAndRun()
}

func refreshWindow(w fyne.Window) {
	w.SetTitle(i18n.T("app_title"))
	w.SetMainMenu(nil)
	
	navData := []struct {
		Label string
		Icon  fyne.Resource
	}{
		{i18n.T("tab_embed"), theme.DocumentCreateIcon()},
		{i18n.T("tab_extract"), theme.VisibilityIcon()},
	}
	
	embedPage := NewEmbedTab(w).Content
	extractPage := NewExtractTab(w).Content
	
	pages := []fyne.CanvasObject{embedPage, extractPage}
	
	contentContainer := container.NewStack()
	
	navList := widget.NewList(
		func() int { return len(navData) },
		func() fyne.CanvasObject {
			
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel("Template Label"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			box := item.(*fyne.Container)
			icon := box.Objects[0].(*widget.Icon)
			label := box.Objects[1].(*widget.Label)
			
			icon.SetResource(navData[id].Icon)
			label.SetText(navData[id].Label)
			
			label.TextStyle = fyne.TextStyle{Bold: true}
		},
	)
	
	navList.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(pages) {
			return
		}
		contentContainer.Objects = []fyne.CanvasObject{pages[id]}
		contentContainer.Refresh()
	}
	
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
				if code == "my" {
					fyne.CurrentApp().Settings().SetTheme(&myTheme{})
				} else {
					fyne.CurrentApp().Settings().SetTheme(theme.DefaultTheme())
				}
				refreshWindow(w)
			}
		}
	})
	if label, ok := codeMap[i18n.GetLang()]; ok {
		langSelect.SetSelected(label)
	} else {
		langSelect.SetSelected("Chinese")
	}
	
	var header fyne.CanvasObject
	if data, err := iconData.ReadFile("assets/Icon.png"); err == nil {
		res := fyne.NewStaticResource("Icon.png", data)
		img := canvas.NewImageFromResource(res)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(24, 24))
		
		label := widget.NewLabelWithStyle("ZUON", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Monospace: true})
		
		header = container.NewPadded(container.NewCenter(container.NewHBox(img, label)))
	} else {
		header = container.NewPadded(widget.NewLabelWithStyle("ZUON", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true}))
	}
	
	var langIcon fyne.CanvasObject
	if data, err := iconData.ReadFile("assets/language.png"); err == nil {
		res := fyne.NewStaticResource("language.png", data)
		img := canvas.NewImageFromResource(res)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(20, 20))
		langIcon = img
	} else {
		langIcon = widget.NewIcon(theme.SettingsIcon())
	}
	
	footer := container.NewPadded(
		container.NewVBox(
			widget.NewSeparator(),
			container.NewBorder(nil, nil, langIcon, nil, langSelect),
		),
	)
	
	sidebar := container.NewBorder(
		header,
		footer,
		nil, nil,
		navList,
	)
	
	split := container.NewHSplit(sidebar, contentContainer)
	split.SetOffset(0.2)
	
	w.SetContent(split)
	
	navList.Select(0)
}
