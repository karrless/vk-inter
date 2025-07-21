package utils

import (
	"io"
	"net/http"
	"time"
)

func ValidateImageURL(imageURL string) error {
	const maxContentLength int64 = 50 * 1024 * 1024 // 50 MB
	const maxBytesToRead int64 = 512 * 1024

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(imageURL)
	if err != nil {
		return newInvalidImageURL()
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return newInvalidImageURL()
	}

	// Optional size check
	if resp.ContentLength > 0 && resp.ContentLength > maxContentLength {
		return newImageTooLarge()
	}

	// Read limited bytes to detect content type
	limitedReader := io.LimitReader(resp.Body, maxBytesToRead)
	buffer := make([]byte, 512)
	n, err := limitedReader.Read(buffer)
	if err != nil && err != io.EOF {
		return newInvalidImageURL()
	}

	contentType := http.DetectContentType(buffer[:n])
	if contentType == "" || contentType[:6] != "image/" {
		return newInvalidImageMimeType()
	}

	return nil
}
