package main

import (
	"fmt"
	"mime"
)

// formatBytes converts bytes to human readable string (KB, MB)
func formatBytes(s int) string {
	if s < 1024 {
		return fmt.Sprintf("%d B", s)
	} else if s < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(s)/1024)
	}
	return fmt.Sprintf("%.2f MB", float64(s)/(1024*1024))
}

// suggestExtension guesses the file extension based on Content-Type
func suggestExtension(contentType string) string {
	switch contentType {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	case "application/zip":
		return ".zip"
	case "application/pdf":
		return ".pdf"
	case "text/plain; charset=utf-8":
		return ".txt"
	}

	exts, _ := mime.ExtensionsByType(contentType)
	if len(exts) > 0 {
		return exts[0]
	}
	return ".bin"
}
