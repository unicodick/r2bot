package archiver

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
)

type ZipArchiver struct{}

func NewZipArchiver() *ZipArchiver {
	return &ZipArchiver{}
}

func (z *ZipArchiver) CreateZip(ctx context.Context, files []FileData) (io.Reader, int64, error) {
	if len(files) == 0 {
		return nil, 0, fmt.Errorf("no files provided for archiving")
	}

	buf := new(bytes.Buffer)

	zipWriter := zip.NewWriter(buf)

	for _, file := range files {
		// check if context was cancelled
		select {
		case <-ctx.Done():
			zipWriter.Close()
			return nil, 0, ctx.Err()
		default:
		}

		zipFile, err := zipWriter.Create(file.Name)
		if err != nil {
			zipWriter.Close()
			return nil, 0, fmt.Errorf("failed to create file %s in zip: %w", file.Name, err)
		}

		_, err = io.Copy(zipFile, file.Reader)
		if err != nil {
			zipWriter.Close()
			return nil, 0, fmt.Errorf("failed to write file %s to zip: %w", file.Name, err)
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to finalize zip archive: %w", err)
	}

	size := int64(buf.Len())
	reader := bytes.NewReader(buf.Bytes())

	return reader, size, nil
}
