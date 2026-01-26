package archiver

import (
	"context"
	"io"
)

type FileData struct {
	Name   string
	Reader io.Reader
	Size   int64
}

type Archiver interface {
	CreateZip(ctx context.Context, files []FileData) (io.Reader, int64, error)
}
