package usecase

import (
	"io"
	"log"
)

type DownloadTelegramFile struct {
	urlGetter  TelegramURLGetter
	downloader HTTPDownloader
}

type TelegramURLGetter interface {
	GetFileURL(fileID string) (string, error)
}

type HTTPDownloader interface {
	Download(url string) (io.ReadCloser, int64, error)
}

func NewDownloadTelegramFile(urlGetter TelegramURLGetter, downloader HTTPDownloader) *DownloadTelegramFile {
	return &DownloadTelegramFile{
		urlGetter:  urlGetter,
		downloader: downloader,
	}
}

func (d *DownloadTelegramFile) Execute(fileID string) (io.ReadCloser, int64, error) {
	url, err := d.urlGetter.GetFileURL(fileID)
	if err != nil {
		log.Printf("failed to get file URL: %v", err)
		return nil, 0, err
	}

	reader, size, err := d.downloader.Download(url)
	if err != nil {
		log.Printf("failed to download file: %v", err)
		return nil, 0, err
	}

	return reader, size, nil
}
