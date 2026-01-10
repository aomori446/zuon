package ui

import (
	"embed"
	"net/url"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/aomori446/zuon/front/i18n"
	"github.com/aomori446/zuon/front/ui/core"
	"github.com/aomori446/zuon/front/ui/pages"
)

//go:embed core/assets/*.png
var iconData embed.FS

func Start() {
	a := app.NewWithID("com.aomori446.zuon")
	i18n.Init()
	
	core.ApplyTheme(a)
	
	token := a.Preferences().String("auth_token")
	if token == "" {
		pages.ShowLoginWindow(a, func(newToken string) {
			a.Preferences().SetString("auth_token", newToken)
			showMainWindow(a)
		})
	} else {
		showMainWindow(a)
	}
	
	a.Run()
}

func showMainWindow(a fyne.App) {
	w := a.NewWindow(i18n.T("app_title"))
	refreshWindow(w)
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(900, 600))
	w.Show()
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
		{i18n.T("tab_settings"), theme.SettingsIcon()},
	}
	
	embedPage := pages.NewEmbedTab(w).Content
	extractPage := pages.NewExtractTab(w).Content
	settingsPage := pages.NewSettingsTab(fyne.CurrentApp(), func() {
		refreshWindow(w)
	}, func() {
		// onLogout
		a := fyne.CurrentApp()
		a.Preferences().SetString("auth_token", "")
		w.Close()
		pages.ShowLoginWindow(a, func(token string) {
			a.Preferences().SetString("auth_token", token)
			showMainWindow(a)
		})
	}).Content
	
	pageObjects := []fyne.CanvasObject{embedPage, extractPage, settingsPage}
	
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
		if id < 0 || id >= len(pageObjects) {
			return
		}
		contentContainer.Objects = []fyne.CanvasObject{pageObjects[id]}
		contentContainer.Refresh()
	}
	
	var header fyne.CanvasObject
	if data, err := iconData.ReadFile("core/assets/Icon.png"); err == nil {
		res := fyne.NewStaticResource("Icon.png", data)
		img := canvas.NewImageFromResource(res)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(24, 24))
		
		label := widget.NewLabelWithStyle("ZUON", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Monospace: true})
		
		header = container.NewPadded(container.NewCenter(container.NewHBox(img, label)))
	} else {
		header = container.NewPadded(widget.NewLabelWithStyle("ZUON", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true}))
	}
	
	ghUrl, _ := url.Parse("https://github.com/aomori446/zuon")
	ghLink := widget.NewHyperlink("Star on GitHub", ghUrl)
	ghLink.Alignment = fyne.TextAlignCenter

	footer := container.NewPadded(
		container.NewVBox(
			widget.NewSeparator(),
			ghLink,
			widget.NewLabelWithStyle(core.AppVersion, fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
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
