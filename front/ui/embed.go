package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/aomori446/zuon/internal"
)

// EmbedTab creates the UI for embedding data.
// It returns the container and a reset function.
func EmbedTab(w fyne.Window) (fyne.CanvasObject, func()) {
	var selectedImage image.Image
	var selectedImageName string // Store the original filename
	var fileToHide []byte
	var fileExtension string

	// Forward declaration to allow use in reset
	var selectBtn *widget.Button
	
	selectBtn, _ = createImageSelector(w, func(img image.Image, name string) {
		selectedImage = img
		selectedImageName = name
	})

	passEntry := widget.NewPasswordEntry()
	passEntry.PlaceHolder = "Encryption Password (min 6 chars)"
	passEntry.OnChanged = func(s string) {
		if err := internal.ValidatePassword(s); err != nil {
			passEntry.SetValidationError(err)
		} else {
			passEntry.SetValidationError(nil)
		}
	}

	msgEntry := widget.NewMultiLineEntry()
	msgEntry.PlaceHolder = "Enter secret message here..."
	msgEntry.SetMinRowsVisible(5)

	var selectFileBtn *widget.Button
	selectFileBtn = widget.NewButtonWithIcon("Select File to Hide", theme.FileIcon(), func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			data, err := io.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			fileToHide = data
			fileExtension = filepath.Ext(reader.URI().Name())
			selectFileBtn.SetText(fmt.Sprintf("%s (%s)", reader.URI().Name(), formatBytes(len(data))))
		}, w)
		fd.Show()
	})

	textContainer := container.NewPadded(msgEntry)
	fileContainer := container.NewPadded(selectFileBtn)
	fileContainer.Hide()

	inputType := widget.NewRadioGroup([]string{"Text", "File"}, func(s string) {
		if s == "Text" {
			fileContainer.Hide()
			textContainer.Show()
		} else {
			textContainer.Hide()
			fileContainer.Show()
		}
	})
	inputType.Horizontal = true
	inputType.Selected = "Text"

	embedBtn := widget.NewButtonWithIcon("Embed & Save", theme.DocumentSaveIcon(), func() {
		if selectedImage == nil {
			dialog.ShowError(errors.New("please select an image first"), w)
			return
		}
		if err := internal.ValidatePassword(passEntry.Text); err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Warning Dialog before proceeding
		dialog.ShowConfirm("Important", "The output image MUST be saved as a PNG file to preserve the hidden data.\n\nDo you want to proceed?", func(ok bool) {
			if !ok {
				return
			}

			var data []byte
			var ext string
			
			if inputType.Selected == "Text" {
				if msgEntry.Text == "" {
					dialog.ShowError(errors.New("message is empty"), w)
					return
				}
				data = []byte(msgEntry.Text)
				ext = ".txt" // Default extension for text
			} else {
				if len(fileToHide) == 0 {
					dialog.ShowError(errors.New("no file selected"), w)
					return
				}
				data = fileToHide
				ext = fileExtension
			}

			res, err := internal.EmbedData(selectedImage, data, ext, 0, passEntry.Text)
			if err != nil {
				dialog.ShowError(fmt.Errorf("embed error: %w", err), w)
				return
			}

			sd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if writer == nil {
					return
				}
				defer writer.Close()

				if err := png.Encode(writer, res); err != nil {
					dialog.ShowError(err, w)
					return
				}
				dialog.ShowInformation("Success", "Image saved successfully!", w)
			}, w)
			
			// Generate suggested filename: original_name_zuon.png
			baseName := selectedImageName
			if ext := filepath.Ext(baseName); ext != "" {
				baseName = strings.TrimSuffix(baseName, ext)
			}
			if baseName == "" {
				baseName = "secret_image"
			}
			sd.SetFileName(baseName + "_zuon.png")
			
			sd.SetFilter(storage.NewExtensionFileFilter([]string{".png"}))
			sd.Show()
		}, w)
	})
	embedBtn.Importance = widget.HighImportance

	content := container.NewVBox(
		widget.NewLabelWithStyle("1. Select Source Image", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewPadded(selectBtn),

		widget.NewLabelWithStyle("2. Secret Data", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewPadded(container.NewVBox(passEntry, inputType)),
		container.NewStack(textContainer, fileContainer),

		layout.NewSpacer(),
		embedBtn,
	)
	
	resetFunc := func() {
		selectedImage = nil
		selectedImageName = ""
		fileToHide = nil
		fileExtension = ""
		selectBtn.SetText("Select Image")
		selectBtn.SetIcon(theme.FolderOpenIcon())
		passEntry.SetText("")
		msgEntry.SetText("")
		selectFileBtn.SetText("Select File to Hide")
		inputType.SetSelected("Text")
	}
	
	return container.NewPadded(content), resetFunc
}
