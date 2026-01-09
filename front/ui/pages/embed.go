package pages

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"net/url"
	"os"
	"path"
	"time"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/aomori446/zuon/front/i18n"
	"github.com/aomori446/zuon/front/ui/core"
	"github.com/aomori446/zuon/front/ui/widgets"
	"github.com/aomori446/zuon/internal"
)

func NewEmbedTab(parent fyne.Window) *container.TabItem {
	
	var btnImage *widgets.CarryButton
	var labelCapacity *widget.Label
	var cardImage *widget.Card
	
	cardImage, btnImage, labelCapacity = widgets.NewFileSelector(
		parent,
		i18n.T("embed_carrier_title"),
		i18n.T("embed_carrier_subtitle"),
		i18n.T("dialog_select_carrier"),
		[]string{".png", ".jpg", ".jpeg"},
		func(reader fyne.URIReadCloser) {
			img, _, err := image.Decode(reader)
			if err != nil {
				core.ShowLocalizedError(err, parent)
				return
			}
			
			btnImage.Carry = img
			
			capacity := core.FormatBytes(internal.Capacity(img))
			labelCapacity.SetText(i18n.Tf("label_capacity", map[string]interface{}{"Capacity": capacity}))
			labelCapacity.TextStyle = fyne.TextStyle{Bold: true}
			labelCapacity.Show()
		},
	)
	
	unsplashBtn := widget.NewButtonWithIcon(i18n.T("btn_search_web"), theme.SearchIcon(), func() {
		ShowUnsplashSearch(parent, func(img image.Image, name string) {
			btnImage.Carry = img
			btnImage.SetText(name)
			btnImage.SetIcon(theme.ConfirmIcon())
			
			capacity := core.FormatBytes(internal.Capacity(img))
			labelCapacity.SetText(i18n.Tf("label_capacity", map[string]interface{}{"Capacity": capacity}))
			labelCapacity.TextStyle = fyne.TextStyle{Bold: true}
			labelCapacity.Show()
		})
	})
	
	cardImage.Content = container.NewVBox(
		container.NewGridWithColumns(2, btnImage, unsplashBtn),
		labelCapacity,
	)
	
	textEntry := widget.NewMultiLineEntry()
	textEntry.SetPlaceHolder(i18n.T("placeholder_text"))
	textEntry.SetMinRowsVisible(4)
	
	fileSizeLabel := widget.NewLabel("")
	fileSizeLabel.Hide()
	fileSizeLabel.Alignment = fyne.TextAlignCenter
	fileSizeLabel.TextStyle = fyne.TextStyle{Italic: true}
	
	fileBtn := widgets.NewCarryButton(i18n.T("btn_select_file"), theme.FolderOpenIcon())
	fileBtn.OnTapped = func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader == nil {
				return
			}
			if err != nil {
				core.ShowLocalizedError(err, parent)
				return
			}
			defer reader.Close()
			
			if uri := reader.URI(); uri != nil {
				if parentURI, err := storage.Parent(uri); err == nil {
					fyne.CurrentApp().Preferences().SetString("last_opened_dir", parentURI.String())
				}
			}
			
			var sizeStr string
			if uri := reader.URI(); uri != nil && uri.Scheme() == "file" {
				if info, err := os.Stat(uri.Path()); err == nil {
					sizeStr = core.FormatBytes(int(info.Size()))
				}
			}
			
			if sizeStr != "" {
				fileSizeLabel.SetText(i18n.Tf("label_file_size", map[string]interface{}{"Size": sizeStr}))
				fileSizeLabel.Show()
			} else {
				fileSizeLabel.Hide()
			}
			
			fileBtn.Carry = reader.URI()
			fileBtn.SetText(reader.URI().Name())
			fileBtn.SetIcon(theme.ConfirmIcon())
		}, parent)
		d.SetTitleText(i18n.T("dialog_select_hidden_file"))
		
		if lastDir := fyne.CurrentApp().Preferences().String("last_opened_dir"); lastDir != "" {
			if listURI, err := storage.ParseURI(lastDir); err == nil {
				if listable, err := storage.ListerForURI(listURI); err == nil {
					d.SetLocation(listable)
				}
			}
		}
		d.Show()
	}
	fileBtn.Hide()
	textEntry.Hide()
	
	radioGroup := widget.NewRadioGroup([]string{i18n.T("radio_text"), i18n.T("radio_file")}, func(s string) {
		if s == i18n.T("radio_text") {
			textEntry.Show()
			fileBtn.Hide()
			fileSizeLabel.Hide()
		} else {
			textEntry.Hide()
			fileBtn.Show()
			if fileBtn.Carry != nil {
				fileSizeLabel.Show()
			}
		}
	})
	radioGroup.Horizontal = true
	radioGroup.SetSelected(i18n.T("radio_text"))
	
	cardData := widget.NewCard(i18n.T("card_data_title"), i18n.T("card_data_subtitle"),
		container.NewVBox(radioGroup, textEntry, fileBtn, fileSizeLabel),
	)
	
	cardPassword, entryPassword := widgets.NewPasswordCard()
	
	progressBar := widget.NewProgressBarInfinite()
	progressBar.Hide()
	
	embedButton := widget.NewButtonWithIcon(i18n.T("btn_embed_start"), theme.MailSendIcon(), nil)
	embedButton.Importance = widget.HighImportance
	
	embedButton.OnTapped = func() {
		
		if btnImage.Carry == nil {
			core.ShowLocalizedError(internal.ErrNoCarrier, parent)
			return
		}
		
		var data []byte
		var ext string
		
		if radioGroup.Selected == i18n.T("radio_text") {
			text := textEntry.Text
			if text == "" {
				core.ShowLocalizedError(internal.ErrNoText, parent)
				return
			}
			data = []byte(text)
			ext = ""
		} else {
			if fileBtn.Carry == nil {
				core.ShowLocalizedError(internal.ErrNoFile, parent)
				return
			}
			uri := fileBtn.Carry.(fyne.URI)
			f, err := os.Open(uri.Path())
			if err != nil {
				core.ShowLocalizedError(err, parent)
				return
			}
			defer f.Close()
			
			data, err = io.ReadAll(f)
			if err != nil {
				core.ShowLocalizedError(err, parent)
				return
			}
			ext = path.Ext(uri.Path())
		}
		
		if entryPassword.Validate() != nil {
			core.ShowLocalizedError(internal.ErrPasswordShort, parent)
			return
		}
		password := entryPassword.Text
		
		baseImage := btnImage.Carry.(image.Image)
		
		embedButton.Disable()
		progressBar.Show()
		
		go func() {
			embedImg, err := internal.EmbedData(baseImage, data, ext, 0, password)
			
			fyne.Do(func() {
				embedButton.Enable()
				progressBar.Hide()
				
				if err != nil {
					core.ShowLocalizedError(err, parent)
					return
				}
				
				saveEmbedResult(parent, embedImg)
			})
		}()
	}
	
	contentVBox := container.New(
		&widgets.CustomVBox{},
		cardImage,
		layout.NewSpacer(),
		cardData,
		layout.NewSpacer(),
		cardPassword,
		layout.NewSpacer(),
		progressBar,
		embedButton,
	)
	
	return container.NewTabItemWithIcon(i18n.T("tab_embed"), theme.DocumentCreateIcon(), container.NewScroll(container.NewPadded(contentVBox)))
}

func saveEmbedResult(parent fyne.Window, img image.Image) {
	
	fsDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if writer == nil {
			return
		}
		if err != nil {
			core.ShowLocalizedError(err, parent)
			return
		}
		defer writer.Close()
		
		err = png.Encode(writer, img)
		if err != nil {
			core.ShowLocalizedError(err, parent)
			return
		}
		
		savedTo := writer.URI().Path()
		hyperlink := widget.NewHyperlink(savedTo, &url.URL{
			Scheme: "file",
			Path:   savedTo,
		})
		
		dialog.NewCustom(i18n.T("dialog_embed_success"), "OK",
			container.NewVBox(
				widget.NewLabel(i18n.T("dialog_file_saved_to")),
				hyperlink,
			), parent).Show()
		
	}, parent)
	
	fsDialog.SetTitleText(i18n.T("dialog_save_embed_title"))
	fsDialog.SetFileName(fmt.Sprintf("%d_zuon.png", time.Now().Unix()))
	fsDialog.SetFilter(storage.NewExtensionFileFilter([]string{".png"}))
	fsDialog.Show()
}
