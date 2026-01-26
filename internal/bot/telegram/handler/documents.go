package handler

import (
	"bytes"
	"io"

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

	// retard var buffer for me
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	if err != nil {
		h.api.SendText(message.Chat.ID, "failed to buffer file")
		return
	}

	fileData := archiver.FileData{
		Name:   filename,
		Reader: bytes.NewReader(buf.Bytes()),
		Size:   size,
	}

	h.mediaGroupMgr.AddFile(message, fileData)
}
