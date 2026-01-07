package ui

import (
	"errors"
	"image/png"
	"os"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
		"fyne.io/fyne/v2/theme"
		"fyne.io/fyne/v2/widget"
		"github.com/aomori446/zuon/front/i18n"
		"github.com/aomori446/zuon/internal"
	)
	
	func NewExtractTab(parent fyne.Window) *container.TabItem {
		
		var btnImage *CarryButton
		var cardImage *widget.Card
	
		cardImage, btnImage, _ = NewFileSelector(
			parent,
			i18n.T("extract_source_title"),
			i18n.T("extract_source_subtitle"),
			i18n.T("dialog_select_extract_source"),
			[]string{".png"},
			func(reader fyne.URIReadCloser) {
				btnImage.Carry = reader.URI()
			},
		)
		
		cardPassword, entryPassword := NewPasswordCard()
		
		progressBar := widget.NewProgressBarInfinite()
		progressBar.Hide()
		
		extractButton := widget.NewButtonWithIcon(i18n.T("btn_extract_start"), theme.MediaReplayIcon(), nil)
		extractButton.Importance = widget.HighImportance
		
		extractButton.OnTapped = func() {
			
			if btnImage.Carry == nil {
				dialog.ShowError(errors.New(i18n.T("err_no_extract_source")), parent)
				return
			}
			
			if entryPassword.Validate() != nil {
				dialog.ShowError(errors.New(i18n.T("err_password_short")), parent)
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
						dialog.ShowError(err, parent)
					})
					return
				}
				defer f.Close()
				
				img, err := png.Decode(f)
				if err != nil {
					fyne.Do(func() {
						extractButton.Enable()
						progressBar.Hide()
						dialog.ShowError(err, parent)
					})
					return
				}
				
				data, ext, err := internal.ExtractData(img, 0, password)
				
				fyne.Do(func() {
					extractButton.Enable()
					progressBar.Hide()
					
					if err != nil {
						dialog.ShowError(err, parent)
						return
					}
					
					ShowResultDialog(parent, data, ext)
				})
			}()
		}
		
		contentVBox := container.New(
			&customVBox{},
			cardImage,
			layout.NewSpacer(),
			cardPassword,
			layout.NewSpacer(),
			progressBar,
			extractButton,
		)
		
		return container.NewTabItem(i18n.T("tab_extract"), container.NewScroll(container.NewPadded(contentVBox)))
	}
	