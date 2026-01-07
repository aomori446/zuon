package unsplash

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"time"
)

func DownloadBytes(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image: %d", resp.StatusCode)
	}
	
	return io.ReadAll(resp.Body)
}

func DownloadImage(url string) (image.Image, error) {
	data, err := DownloadBytes(url)
	if err != nil {
		return nil, err
	}
	
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	
	return img, nil
}
