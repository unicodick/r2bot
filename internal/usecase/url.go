package usecase

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"path"
	"strings"

	"github.com/unicodick/r2bot/internal/storage"
	"github.com/unicodick/r2bot/internal/utils"
)

type UploadURL struct {
	uploader   storage.Uploader
	downloader HTTPDownloader
}

func NewUploadURL(uploader storage.Uploader, downloader HTTPDownloader) *UploadURL {
	return &UploadURL{
		uploader:   uploader,
		downloader: downloader,
	}
}

func (u *UploadURL) Execute(ctx context.Context, fileURL string) (string, string, error) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", "", fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	filename := u.extractFilename(parsedURL)
	if filename == "" {
		filename = "downloaded_file"
	}

	// download file with zero-copy streaming
	reader, size, err := u.downloader.Download(fileURL)
	if err != nil {
		log.Printf("download from URL failed: %v", err)
		return "", "", fmt.Errorf("download failed: %w", err)
	}
	defer reader.Close()

	if size > 50*1024*1024 {
		return "", "", fmt.Errorf("file too big (>50MB)")
	}

	fileInfo, err := u.uploader.UploadFile(ctx, filename, reader, size)
	if err != nil {
		log.Printf("upload to R2 failed: %v", err)
		return "", "", fmt.Errorf("upload failed: %w", err)
	}

	text := "file: " + fileInfo.Name + "\nsize: " + utils.FormatFileSize(fileInfo.Size) + "\nlink: " + fileInfo.URL
	return text, fileInfo.URL, nil
}

func (u *UploadURL) extractFilename(parsedURL *url.URL) string {
	filename := path.Base(parsedURL.Path)

	// rm query parameters and fragments
	if idx := strings.Index(filename, "?"); idx != -1 {
		filename = filename[:idx]
	}
	if idx := strings.Index(filename, "#"); idx != -1 {
		filename = filename[:idx]
	}

	if filename == "." || filename == "/" || filename == "" {
		return ""
	}

	return filename
}
