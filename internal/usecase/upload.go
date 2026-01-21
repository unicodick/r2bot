package usecase

import (
	"context"
	"io"
	"log"

	"github.com/unicodick/r2bot/internal/storage"
	"github.com/unicodick/r2bot/internal/utils"
)

type UploadFile struct {
	uploader storage.Uploader
}

func NewUploadFile(uploader storage.Uploader) *UploadFile {
	return &UploadFile{
		uploader: uploader,
	}
}

func (u *UploadFile) Execute(ctx context.Context, filename string, reader io.Reader, size int64) (string, string, error) {
	fileInfo, err := u.uploader.UploadFile(ctx, filename, reader, size)
	if err != nil {
		log.Printf("upload failed: %v", err)
		return "", "", err
	}

	text := "file: " + fileInfo.Name + "\nsize: " + utils.FormatFileSize(fileInfo.Size) + "\nlink: " + fileInfo.URL
	return text, fileInfo.URL, nil
}
