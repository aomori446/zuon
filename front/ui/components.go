package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/aomori446/zuon/internal"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

func createImageSelector(parent fyne.Window, onImageLoaded func(image.Image, string)) (*widget.Button, *widget.Label) {
	var btn *widget.Button
	btn = widget.NewButtonWithIcon("Select Image", theme.FolderOpenIcon(), func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()

			img, _, err := image.Decode(reader)
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}

			capacity := internal.MaxCapacity(img)
			btn.SetText(fmt.Sprintf("%s (Max: %s)", reader.URI().Name(), formatBytes(capacity)))

			onImageLoaded(img, reader.URI().Name())
		}, parent)

		fd.SetFilter(storage.NewExtensionFileFilter([]string{
			".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp", ".tiff", ".tif",
		}))
		fd.Show()
	})
	return btn, nil
}
