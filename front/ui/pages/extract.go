package pages

import (
	"image/png"
	"os"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/aomori446/zuon/front/i18n"
	"github.com/aomori446/zuon/front/ui/core"
	"github.com/aomori446/zuon/front/ui/widgets"
	"github.com/aomori446/zuon/internal"
)

func NewExtractTab(parent fyne.Window) *container.TabItem {
	var btnImage *widgets.CarryButton
	var cardImage *widget.Card
	
	cardImage, btnImage, _ = widgets.NewFileSelector(
		parent,
		i18n.T("extract_source_title"),
		i18n.T("extract_source_subtitle"),
		i18n.T("dialog_select_extract_source"),
		[]string{".png"},
		func(reader fyne.URIReadCloser) {
			btnImage.Carry = reader.URI()
		},
	)
	
	cardPassword, entryPassword := widgets.NewPasswordCard()
	
	progressBar := widget.NewProgressBarInfinite()
	progressBar.Hide()
	
	extractButton := widget.NewButtonWithIcon(i18n.T("btn_extract_start"), theme.MediaReplayIcon(), nil)
	extractButton.Importance = widget.HighImportance
	
	extractButton.OnTapped = func() {
		if btnImage.Carry == nil {
			core.ShowLocalizedError(internal.ErrNoSource, parent)
			return
		}
		
		if entryPassword.Validate() != nil {
			core.ShowLocalizedError(internal.ErrPasswordShort, parent)
			return
		}
		
		uri := btnImage.Carry.(fyne.URI)
		password := entryPassword.Text
		
		extractButton.Disable()
		progressBar.Show()
		
		go func() {
			f, err := os.Open(uri.Path())
			if err != nil {
				fyne.Do(func() {
					extractButton.Enable()
					progressBar.Hide()
					core.ShowLocalizedError(err, parent)
				})
				return
			}
			defer f.Close()
			
			img, err := png.Decode(f)
			if err != nil {
				fyne.Do(func() {
					extractButton.Enable()
					progressBar.Hide()
					core.ShowLocalizedError(err, parent)
				})
				return
			}
			
			data, ext, err := internal.ExtractData(img, 0, password)
			
			fyne.Do(func() {
				extractButton.Enable()
				progressBar.Hide()
				
				if err != nil {
					core.ShowLocalizedError(err, parent)
					return
				}
				
				widgets.ShowResultDialog(parent, data, ext)
			})
		}()
	}
	
	contentVBox := container.New(
		&widgets.CustomVBox{},
		cardImage,
		layout.NewSpacer(),
		cardPassword,
		layout.NewSpacer(),
		progressBar,
		extractButton,
	)
	
	return container.NewTabItemWithIcon(i18n.T("tab_extract"), theme.VisibilityIcon(), container.NewScroll(container.NewPadded(contentVBox)))
}
