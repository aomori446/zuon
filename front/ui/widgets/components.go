package widgets

import (
	"fmt"
	"time"
	
	"bytes"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/aomori446/zuon/front/i18n"
	"github.com/aomori446/zuon/front/ui/core"
	"github.com/aomori446/zuon/internal"
)

type CarryButton struct {
	widget.Button
	Carry interface{}
}

func NewCarryButton(text string, icon fyne.Resource) *CarryButton {
	b := &CarryButton{}
	b.ExtendBaseWidget(b)
	b.SetText(text)
	b.SetIcon(icon)
	return b
}

func NewFileSelector(
	parent fyne.Window,
	title string,
	subtitle string,
	dialogTitle string,
	extensions []string,
	onSelected func(fyne.URIReadCloser),
) (*widget.Card, *CarryButton, *widget.Label) {
	infoLabel := widget.NewLabel("")
	infoLabel.Hide()
	infoLabel.Alignment = fyne.TextAlignCenter
	infoLabel.TextStyle = fyne.TextStyle{Italic: true}
	
	btn := NewCarryButton(i18n.T("btn_select_file"), theme.FolderOpenIcon())
	btn.OnTapped = func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader == nil {
				return
			}
			defer reader.Close()
			if err != nil {
				core.ShowLocalizedError(err, parent)
				return
			}
			
			if uri := reader.URI(); uri != nil {
				parentURI, err := storage.Parent(uri)
				if err == nil {
					fyne.CurrentApp().Preferences().SetString("last_opened_dir", parentURI.String())
				}
			}
			
			btn.SetText(reader.URI().Name())
			btn.SetIcon(theme.ConfirmIcon())
			
			if onSelected != nil {
				onSelected(reader)
			}
		}, parent)
		
		d.SetTitleText(dialogTitle)
		if len(extensions) > 0 {
			d.SetFilter(storage.NewExtensionFileFilter(extensions))
		}
		
		lastDir := fyne.CurrentApp().Preferences().String("last_opened_dir")
		if lastDir != "" {
			if listURI, err := storage.ParseURI(lastDir); err == nil {
				if listable, err := storage.ListerForURI(listURI); err == nil {
					d.SetLocation(listable)
				}
			}
		}
		
		d.Show()
	}
	
	content := container.NewVBox(btn, infoLabel)
	card := widget.NewCard(title, subtitle, content)
	return card, btn, infoLabel
}

func NewPasswordCard() (*widget.Card, *widget.Entry) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(i18n.T("placeholder_password"))
	entry.Password = true
	entry.Validator = internal.ValidatePassword
	
	card := widget.NewCard(i18n.T("card_security_title"), i18n.T("card_security_subtitle"), entry)
	return card, entry
}

func ShowResultDialog(parent fyne.Window, data []byte, ext string) {
	if ext == "" {
		entry := widget.NewMultiLineEntry()
		entry.SetText(string(data))
		entry.Wrapping = fyne.TextWrapWord
		
		copyBtn := widget.NewButtonWithIcon(i18n.T("btn_copy_text"), theme.ContentCopyIcon(), func() {
			parent.Clipboard().SetContent(string(data))
		})
		
		content := container.NewBorder(nil, copyBtn, nil, nil, entry)
		custom := dialog.NewCustom(i18n.T("result_title_text"), i18n.T("btn_close"), content, parent)
		custom.Resize(fyne.NewSize(400, 300))
		custom.Show()
		
	} else if ext == ".png" || ext == ".jpeg" || ext == ".jpg" {
		imgContent := canvas.NewImageFromReader(bytes.NewReader(data), fmt.Sprintf("ext_%d", time.Now().Unix()))
		imgContent.FillMode = canvas.ImageFillContain
		imgContent.SetMinSize(fyne.NewSize(300, 300))
		
		saveBtn := widget.NewButtonWithIcon(i18n.T("btn_save_image"), theme.DocumentSaveIcon(), func() {
			saveFile(parent, data, ext)
		})
		
		content := container.NewBorder(nil, saveBtn, nil, nil, imgContent)
		custom := dialog.NewCustom(i18n.T("result_title_image"), i18n.T("btn_close"), content, parent)
		custom.Resize(fyne.NewSize(400, 400))
		custom.Show()
		
	} else {
		info := widget.NewLabel(i18n.Tf("label_format_size", map[string]interface{}{
			"Format": ext,
			"Size":   core.FormatBytes(len(data)),
		}))
		info.Alignment = fyne.TextAlignCenter
		
		saveBtn := widget.NewButtonWithIcon(i18n.T("btn_save_file"), theme.DocumentSaveIcon(), func() {
			saveFile(parent, data, ext)
		})
		
		content := container.NewVBox(info, saveBtn)
		dialog.NewCustom(i18n.T("result_title_file"), i18n.T("btn_close"), content, parent).Show()
	}
}

func saveFile(parent fyne.Window, data []byte, ext string) {
	fsDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if writer == nil {
			return
		}
		if err != nil {
			core.ShowLocalizedError(err, parent)
			return
		}
		defer writer.Close()
		
		_, err = writer.Write(data)
		if err != nil {
			core.ShowLocalizedError(err, parent)
			return
		}
		
		dialog.ShowInformation(i18n.T("dialog_save_success_title"), i18n.T("dialog_file_saved_to")+"\n"+writer.URI().Path(), parent)
	}, parent)
	
	fsDialog.SetFileName(fmt.Sprintf("%d_extracted%s", time.Now().Unix(), ext))
	fsDialog.Show()
}
