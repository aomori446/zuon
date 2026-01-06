package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func main() {
	a := app.NewWithID("com.aomori446.zuon")
	w := a.NewWindow("zuon (ズオン)")

	embedContent, embedReset := EmbedTab(w)
	extractContent, extractReset := ExtractTab(w)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Embed", theme.MailAttachmentIcon(), embedContent),
		container.NewTabItemWithIcon("Extract", theme.VisibilityIcon(), extractContent),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	tabs.OnSelected = func(item *container.TabItem) {
		if item.Text == "Embed" {
			embedReset()
		} else if item.Text == "Extract" {
			extractReset()
		}
	}

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(450, 700))
	w.CenterOnScreen()
	w.ShowAndRun()
}
