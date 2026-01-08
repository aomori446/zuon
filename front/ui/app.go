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

//go:embed assets/Icon.png
var iconData embed.FS

func Start() {
	a := app.NewWithID("com.aomori446.zuon")
	i18n.Init()
	
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
	
	contentContainer := container.NewMax()
	
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
	
	footer := container.NewPadded(
		container.NewVBox(
			widget.NewSeparator(),
			container.NewBorder(nil, nil, widget.NewIcon(theme.SettingsIcon()), nil, langSelect),
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
