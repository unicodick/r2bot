package services

import (
	"fmt"
	"io"

	"github.com/unicodick/r2bot/internal/http"
)

type FileDownloader struct{}

func NewFileDownloader() *FileDownloader {
	return &FileDownloader{}
}

func (d *FileDownloader) Download(url string) (io.ReadCloser, int64, error) {
	client := http.Client()
	resp, err := client.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("download failed: %w", err)
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, 0, fmt.Errorf("status %d", resp.StatusCode)
	}

	return resp.Body, resp.ContentLength, nil
}
