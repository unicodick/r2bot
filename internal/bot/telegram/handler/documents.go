package handler

import (
	"io"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/unicodick/r2bot/internal/archiver"
)

func (h *Handler) handleDocument(message *tgbotapi.Message) {
	if message.Document.FileSize > 50*1024*1024 {
		h.api.SendText(message.Chat.ID, "file too big (>50MB)")
		return
	}

	reader, size, err := h.downloadTg.Execute(message.Document.FileID)
	if err != nil {
		h.api.SendText(message.Chat.ID, "failed to download file")
		return
	}
	defer reader.Close()

	filename := h.getFileName(message)

	tempFile, err := os.CreateTemp("", "telegram_file_*.tmp")
	if err != nil {
		h.api.SendText(message.Chat.ID, "failed to create temporary file")
		return
	}

	_, err = io.Copy(tempFile, reader)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		h.api.SendText(message.Chat.ID, "failed to write temporary file")
		return
	}

	_, err = tempFile.Seek(0, 0)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		h.api.SendText(message.Chat.ID, "failed to reset file position")
		return
	}

	fileData := archiver.FileData{
		Name:   filename,
		Reader: &tempFileReader{file: tempFile, path: tempFile.Name()},
		Size:   size,
	}

	h.mediaGroupMgr.AddFile(message, fileData)
}

type tempFileReader struct {
	file *os.File
	path string
}

func (t *tempFileReader) Read(p []byte) (int, error) {
	return t.file.Read(p)
}

func (t *tempFileReader) Close() error {
	err := t.file.Close()
	os.Remove(t.path)
	return err
}
