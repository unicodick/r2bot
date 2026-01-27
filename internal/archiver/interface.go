package archiver

import (
	"context"
	"io"
)

type FileData struct {
	Name   string
	Reader io.ReadCloser
	Size   int64
}

type Archiver interface {
	CreateZip(ctx context.Context, files []FileData) (io.Reader, int64, error)
}

// wraps an io.Reader to implement io.ReadCloser with a no-op close method
type NopCloser struct {
	io.Reader
}

func (NopCloser) Close() error {
	return nil
}

func WrapReader(r io.Reader) io.ReadCloser {
	if rc, ok := r.(io.ReadCloser); ok {
		return rc
	}
	return NopCloser{r}
}
