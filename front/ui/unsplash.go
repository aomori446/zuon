package ui

import (
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/aomori446/zuon/front/i18n"
	"github.com/aomori446/zuon/internal"
	"github.com/aomori446/zuon/internal/unsplash"
)

type tappableImage struct {
	widget.BaseWidget
	img      *canvas.Image
	onTapped func()
}

func newTappableImageFromImage(img image.Image, onTapped func()) *tappableImage {
	t := &tappableImage{
		img:      canvas.NewImageFromImage(img),
		onTapped: onTapped,
	}
	t.img.FillMode = canvas.ImageFillContain
	t.img.SetMinSize(fyne.NewSize(150, 100))
	t.ExtendBaseWidget(t)
	return t
}

func (t *tappableImage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.img)
}

func (t *tappableImage) Tapped(_ *fyne.PointEvent) {
	if t.onTapped != nil {
		t.onTapped()
	}
}

func ShowUnsplashSearch(parent fyne.Window, onSelected func(img image.Image, name string)) {
	
	apiKey := os.Getenv("UNSPLASH_ACCESS_KEY")
	if apiKey == "" {
		apiKey = fyne.CurrentApp().Preferences().String("unsplash_key")
	}
	
	if apiKey == "" {
		input := widget.NewEntry()
		input.SetPlaceHolder("Access Key")
		
		msgLabel := widget.NewLabel(i18n.T("dialog_api_key_msg"))
		msgLabel.Wrapping = fyne.TextWrapWord
		
		loading := widget.NewProgressBarInfinite()
		loading.Hide()
		
		var d dialog.Dialog
		
		confirmBtn := widget.NewButton(i18n.T("btn_save"), func() {
			key := input.Text
			if key == "" {
				return
			}
			
			input.Disable()
			loading.Show()
			
			go func() {
				
				client, err := unsplash.NewClient(key)
				if err == nil {
					
					_, err = client.SearchPhotos("test", 1, 1)
				}
				
				fyne.Do(func() {
					input.Enable()
					loading.Hide()
					
					if err != nil {
						ShowLocalizedError(err, parent)
					} else {
						
						fyne.CurrentApp().Preferences().SetString("unsplash_key", key)
						d.Hide()
						
						ShowUnsplashSearch(parent, onSelected)
					}
				})
			}()
		})
		confirmBtn.Importance = widget.HighImportance
		
		content := container.NewVBox(msgLabel, input, loading, confirmBtn)
		
		d = dialog.NewCustom(i18n.T("dialog_api_key_title"), i18n.T("btn_close"), content, parent)
		d.Resize(fyne.NewSize(400, 200))
		d.Show()
		return
	}
	
	client, err := unsplash.NewClient(apiKey)
	if err != nil {
		ShowLocalizedError(err, parent)
		return
	}
	
	var d dialog.Dialog
	
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder(i18n.T("unsplash_placeholder"))
	
	resultsContainer := container.NewGridWrap(fyne.NewSize(180, 150))
	scrollContainer := container.NewVScroll(resultsContainer)
	scrollContainer.SetMinSize(fyne.NewSize(600, 400))
	
	loading := widget.NewProgressBarInfinite()
	loading.Hide()
	
	searchBtn := widget.NewButtonWithIcon(i18n.T("btn_search"), theme.SearchIcon(), nil)
	
	performSearch := func() {
		query := searchEntry.Text
		if query == "" {
			return
		}
		
		searchBtn.Disable()
		loading.Show()
		resultsContainer.Objects = nil
		resultsContainer.Refresh()
		
		go func() {
			defer func() {
				fyne.Do(func() {
					searchBtn.Enable()
					loading.Hide()
				})
			}()
			
			results, err := client.SearchPhotos(query, 1, 12)
			
			if err != nil {
				fyne.Do(func() {
					
					if errors.Is(err, internal.ErrInvalidAPIKey) {
						fyne.CurrentApp().Preferences().SetString("unsplash_key", "")
						if d != nil {
							d.Hide()
						}
					}
					ShowLocalizedError(err, parent)
				})
				return
			}
			
			for _, photo := range results.Results {
				photo := photo
				
				fyne.Do(func() {
					thumbContainer := container.NewVBox()
					authorLabel := widget.NewLabel(fmt.Sprintf("@%s", photo.User.Username))
					authorLabel.Alignment = fyne.TextAlignCenter
					authorLabel.Truncation = fyne.TextTruncateEllipsis
					
					onImageTapped := func() {
						waitDialog := dialog.NewCustomWithoutButtons(i18n.T("tab_embed"), widget.NewProgressBarInfinite(), parent)
						waitDialog.Show()
						
						go func() {
							fullImg, err := unsplash.DownloadImage(photo.URLs.Regular)
							
							fyne.Do(func() {
								waitDialog.Hide()
								if err != nil {
									ShowLocalizedError(err, parent)
									return
								}
								parent.Canvas().Refresh(parent.Content())
								
								onSelected(fullImg, photo.ID+".jpg")
								if d != nil {
									d.Hide()
								}
							})
						}()
					}
					
					loadingContainer := container.NewCenter(widget.NewProgressBarInfinite())
					
					thumbContainer.Add(loadingContainer)
					thumbContainer.Add(authorLabel)
					resultsContainer.Add(thumbContainer)
					resultsContainer.Refresh()
					
					go func(container *fyne.Container, p unsplash.Photo) {
						img, err := unsplash.DownloadImage(p.URLs.Thumb)
						if err != nil {
							log.Printf("Failed to download thumb: %v", err)
							return
						}
						
						fyne.Do(func() {
							
							tImg := newTappableImageFromImage(img, onImageTapped)
							
							if len(container.Objects) > 0 {
								container.Objects[0] = tImg
								container.Refresh()
							}
						})
					}(thumbContainer, photo)
				})
			}
		}()
	}
	
	searchBtn.OnTapped = performSearch
	searchEntry.OnSubmitted = func(s string) { performSearch() }
	
	top := container.NewBorder(nil, nil, nil, searchBtn, searchEntry)
	content := container.NewBorder(container.NewVBox(top, loading), nil, nil, nil, scrollContainer)
	
	d = dialog.NewCustom(i18n.T("unsplash_title"), i18n.T("btn_close"), content, parent)
	d.Resize(fyne.NewSize(650, 550))
	d.Show()
}
