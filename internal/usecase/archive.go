package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/unicodick/r2bot/internal/archiver"
	"github.com/unicodick/r2bot/internal/storage"
	"github.com/unicodick/r2bot/internal/utils"
)

type UploadArchive struct {
	uploader storage.Uploader
	archiver archiver.Archiver
}

func NewUploadArchive(uploader storage.Uploader, archiver archiver.Archiver) *UploadArchive {
	return &UploadArchive{
		uploader: uploader,
		archiver: archiver,
	}
}

func (u *UploadArchive) Execute(ctx context.Context, files []archiver.FileData, archiveName string) (string, string, error) {
	if len(files) == 0 {
		return "", "", fmt.Errorf("no files provided for archiving")
	}

	if archiveName == "" {
		archiveName = fmt.Sprintf("archive_%d.zip", time.Now().Unix())
	}

	if len(archiveName) < 4 || archiveName[len(archiveName)-4:] != ".zip" {
		archiveName += ".zip"
	}

	archiveReader, archiveSize, err := u.archiver.CreateZip(ctx, files)
	if err != nil {
		log.Printf("failed to create archive: %v", err)
		return "", "", fmt.Errorf("failed to create archive: %w", err)
	}

	fileInfo, err := u.uploader.UploadFile(ctx, archiveName, archiveReader, archiveSize)
	if err != nil {
		log.Printf("upload failed: %v", err)
		return "", "", err
	}

	filesCount := len(files)
	text := fmt.Sprintf("archive: %s\nfiles: %d\nsize: %s\nlink: %s",
		fileInfo.Name,
		filesCount,
		utils.FormatFileSize(fileInfo.Size),
		fileInfo.URL)

	return text, fileInfo.URL, nil
}
