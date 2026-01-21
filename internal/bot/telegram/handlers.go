package telegram

import (
	"context"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/unicodick/r2bot/internal/usecase"
)

type Handler struct {
	api        *API
	auth       *usecase.CheckAuth
	downloadTg *usecase.DownloadTelegramFile
	upload     *usecase.UploadFile
}

func NewHandler(api *API, auth *usecase.CheckAuth, downloadTg *usecase.DownloadTelegramFile, upload *usecase.UploadFile) *Handler {
	return &Handler{
		api:        api,
		auth:       auth,
		downloadTg: downloadTg,
		upload:     upload,
	}
}

func (h *Handler) HandleMessage(message *tgbotapi.Message) {
	if !h.auth.Execute(message.From.ID) {
		return
	}

	if message.IsCommand() {
		h.handleCommand(message)
		return
	}

	if message.Document != nil {
		h.handleDocument(message)
	}
}

func (h *Handler) handleCommand(message *tgbotapi.Message) {
	if message.Command() == "start" {
		h.api.SendText(message.Chat.ID, "any uploaded file will be stored in r2")
	}
}

func (h *Handler) getFileName(message *tgbotapi.Message) string {
	if message.Caption != "" {
		caption := strings.TrimSpace(message.Caption)

		if filepath.Ext(caption) == "" {
			originalExt := filepath.Ext(message.Document.FileName)
			if originalExt != "" {
				caption += originalExt
			}
		}

		return caption
	}

	// fallback
	return message.Document.FileName
}

func (h *Handler) handleDocument(message *tgbotapi.Message) {
	if message.Document.FileSize > 50*1024*1024 {
		h.api.SendText(message.Chat.ID, "file too big (<50mb)")
		return
	}

	reader, size, err := h.downloadTg.Execute(message.Document.FileID)
	if err != nil {
		h.api.SendText(message.Chat.ID, "failed to download file")
		return
	}
	defer reader.Close()

	filename := h.getFileName(message)
	text, url, err := h.upload.Execute(context.Background(), filename, reader, size)
	if err != nil {
		h.api.SendText(message.Chat.ID, "failed to upload file")
		return
	}

	h.api.SendWithButton(message.Chat.ID, text, "✈️", url)

	if err := h.api.DeleteMessage(message.Chat.ID, message.MessageID); err != nil {
		// intentionally ignored - message deletion is optional cleanup
	}
}
