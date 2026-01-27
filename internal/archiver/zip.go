package archiver

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"sync"
)

type ZipArchiver struct{}

func NewZipArchiver() *ZipArchiver {
	return &ZipArchiver{}
}

func (z *ZipArchiver) CreateZip(ctx context.Context, files []FileData) (io.Reader, int64, error) {
	if len(files) == 0 {
		return nil, 0, fmt.Errorf("no files provided for archiving")
	}

	pipeReader, pipeWriter := io.Pipe()

	var wg sync.WaitGroup
	var zipError error

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer pipeWriter.Close()

		zipWriter := zip.NewWriter(pipeWriter)

		// track used filenames to avoid duplicates
		usedNames := make(map[string]int)

		for _, file := range files {
			// check if context was cancelled
			select {
			case <-ctx.Done():
				zipWriter.Close()
				zipError = ctx.Err()
				return
			default:
			}

			filename := z.getUniqueFilename(file.Name, usedNames)

			zipFile, err := zipWriter.Create(filename)
			if err != nil {
				zipWriter.Close()
				zipError = fmt.Errorf("failed to create file %s in zip: %w", filename, err)
				return
			}

			_, err = io.Copy(zipFile, file.Reader)
			if closeErr := file.Reader.Close(); closeErr != nil {
				// log the close error but don't fail the entire operation
				// the copy error is more critical
			}
			if err != nil {
				zipWriter.Close()
				zipError = fmt.Errorf("failed to write file %s to zip: %w", filename, err)
				return
			}
		}

		err := zipWriter.Close()
		if err != nil {
			zipError = fmt.Errorf("failed to finalize zip archive: %w", err)
			return
		}
	}()

	streamReader := &streamingZipReader{
		reader: pipeReader,
		wg:     &wg,
		err:    &zipError,
	}

	return streamReader, -1, nil
}

// wraps the pipe reader and handles goroutine synchronization
type streamingZipReader struct {
	reader io.ReadCloser
	wg     *sync.WaitGroup
	err    *error
	once   sync.Once
}

func (s *streamingZipReader) Read(p []byte) (int, error) {
	n, err := s.reader.Read(p)

	if err == io.EOF {
		s.once.Do(func() {
			s.wg.Wait()
		})

		if *s.err != nil {
			return n, *s.err
		}
	}

	return n, err
}

func (s *streamingZipReader) Close() error {
	s.once.Do(func() {
		s.wg.Wait()
	})
	return s.reader.Close()
}

func (z *ZipArchiver) getUniqueFilename(originalName string, usedNames map[string]int) string {
	if _, exists := usedNames[originalName]; !exists {
		usedNames[originalName] = 1
		return originalName
	}

	// if the name is used, find a unique variant
	ext := filepath.Ext(originalName)
	nameWithoutExt := originalName[:len(originalName)-len(ext)]

	counter := usedNames[originalName] + 1

	for {
		newName := fmt.Sprintf("%s_%s%s", nameWithoutExt, strconv.Itoa(counter), ext)
		if _, exists := usedNames[newName]; !exists {
			usedNames[originalName] = counter
			usedNames[newName] = 1
			return newName
		}
		counter++
	}
}
