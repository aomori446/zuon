package main

import (
	"errors"
	"fmt"
	"image"
	"net/http"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/aomori446/zuon/internal"
)

// ExtractTab creates the UI for extracting data.
// It returns the container and a reset function.
func ExtractTab(w fyne.Window) (fyne.CanvasObject, func()) {
	var selectedImage image.Image
	var extractedData []byte
	var extractedExt string

	// Forward declaration
	var selectBtn *widget.Button
	
	selectBtn, _ = createImageSelector(w, func(img image.Image, name string) {
		selectedImage = img
	})

	passEntry := widget.NewPasswordEntry()
	passEntry.PlaceHolder = "Decryption Password"
	passEntry.OnChanged = func(s string) {
		if err := internal.ValidatePassword(s); err != nil {
			passEntry.SetValidationError(err)
		} else {
			passEntry.SetValidationError(nil)
		}
	}

	resultCard := widget.NewCard("Extraction Result", "", nil)
	resultCard.Hide()

	extractBtn := widget.NewButtonWithIcon("Extract Data", theme.ConfirmIcon(), func() {
		if selectedImage == nil {
			dialog.ShowError(errors.New("please select an image first"), w)
			return
		}
		if err := internal.ValidatePassword(passEntry.Text); err != nil {
			dialog.ShowError(err, w)
			return
		}

		data, ext, err := internal.ExtractData(selectedImage, 0, passEntry.Text)
		if err != nil {
			resultCard.Hide()
			if errors.Is(err, internal.ErrDecryptionFailed) {
				dialog.ShowError(errors.New("wrong password or no data found"), w)
			} else {
				dialog.ShowError(fmt.Errorf("extract error: %w", err), w)
			}
			return
		}

		extractedData = data
		extractedExt = ext
		
		// Show success dialog? Or just show the result card?
		// Usually showing the result card is enough feedback for "success", 
		// but let's show a small info dialog if preferred, or just rely on the card appearing.
		// The user asked to replace label with dialog for errors/success.
		dialog.ShowInformation("Success", "Data extracted successfully!", w)

		// Analyze Data
		contentType := http.DetectContentType(data)
		isText := utf8.Valid(data) && (contentType == "text/plain; charset=utf-8" || contentType == "text/plain")
		
		// Use stored extension if available, otherwise guess
		suggestedName := "extracted"
		if extractedExt != "" {
			suggestedName += extractedExt
		} else {
			suggestedName += suggestExtension(contentType)
		}

		// Build Result Content
		resultContent := container.NewVBox()

		// Info Row - Simplified: Size and Original Extension (if available)
		infoText := fmt.Sprintf("Size: %s", formatBytes(len(data)))
		if extractedExt != "" {
			infoText += fmt.Sprintf("\nOriginal Extension: %s", extractedExt)
		}
		resultContent.Add(widget.NewLabel(infoText))

		// Action Buttons
		saveBtn := widget.NewButtonWithIcon("Save as File...", theme.DocumentSaveIcon(), func() {
			sd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil || writer == nil {
					return
				}
				defer writer.Close()
				writer.Write(extractedData)
			}, w)
			sd.SetFileName(suggestedName)
			sd.Show()
		})

		btns := container.NewHBox(saveBtn)

		if isText {
			copyBtn := widget.NewButtonWithIcon("Copy Text", theme.ContentCopyIcon(), func() {
				w.Clipboard().SetContent(string(extractedData))
				dialog.ShowInformation("Copied", "Text copied to clipboard", w)
			})
			btns.Add(copyBtn)
		}
		resultContent.Add(btns)

		// Preview
		if isText {
			resultContent.Add(widget.NewSeparator())
			resultContent.Add(widget.NewLabelWithStyle("Text Preview:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))

			textStr := string(data)
			if len(textStr) > 500 {
				textStr = textStr[:500] + "... (truncated)"
			}

			entry := widget.NewMultiLineEntry()
			entry.SetText(textStr)
			entry.Disable()
			entry.TextStyle = fyne.TextStyle{Monospace: true}
			resultContent.Add(entry)
		} else {
			resultContent.Add(widget.NewSeparator())
			resultContent.Add(widget.NewLabelWithStyle("Binary file detected.", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}))
		}

		resultCard.SetContent(resultContent)
		resultCard.Show()
	})
	extractBtn.Importance = widget.HighImportance

	content := container.NewVBox(
		widget.NewLabelWithStyle("1. Select Encoded Image", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewPadded(selectBtn),

		widget.NewLabelWithStyle("2. Decrypt", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewPadded(passEntry),

		extractBtn,
		layout.NewSpacer(),
		resultCard,
	)

	resetFunc := func() {
		selectedImage = nil
		extractedData = nil
		extractedExt = ""
		selectBtn.SetText("Select Image")
		selectBtn.SetIcon(theme.FolderOpenIcon())
		passEntry.SetText("")
		resultCard.Hide()
	}

	return container.NewPadded(container.NewVScroll(content)), resetFunc
}
