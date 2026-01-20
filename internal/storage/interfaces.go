package storage

import "io"

type FileInfo struct {
	Name string
	Size int64
	URL  string
}

type Uploader interface {
	UploadFile(filename string, reader io.Reader, size int64) (*FileInfo, error)
}

type Downloader interface {
	DownloadFile(url string) ([]byte, error)
}
