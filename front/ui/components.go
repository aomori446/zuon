package main

import (
	"fmt"
	"image"
	_ "image/gif" // Register GIF decoder
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/aomori446/zuon/internal"
	_ "golang.org/x/image/bmp"  // Register BMP decoder
	_ "golang.org/x/image/tiff" // Register TIFF decoder
	_ "golang.org/x/image/webp" // Register WebP decoder
)

// createImageSelector creates a button that opens a file dialog to select an image.
// It updates the button text with the filename and capacity.
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
		
		// Expanded filter list
		// Note: Go's image.Decode needs appropriate packages imported for side-effects to support these formats.
		// Standard library supports: png, jpeg, gif.
		// Extended support (requires golang.org/x/image): bmp, tiff, webp.
		// We will add standard ones first, and if you want more, we need to import x/image.
		// For now, let's stick to standard + common ones.
		fd.SetFilter(storage.NewExtensionFileFilter([]string{
			".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp", ".tiff", ".tif",
		}))
		fd.Show()
	})
	return btn, nil
}
