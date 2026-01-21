package telegram

import (
	"context"
	"path/filepath"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/unicodick/r2bot/internal/usecase"
)

type Handler struct {
	api        *API
	auth       *usecase.CheckAuth
	downloadTg *usecase.DownloadTelegramFile
	upload     *usecase.UploadFile
	uploadURL  *usecase.UploadURL
}

func NewHandler(api *API, auth *usecase.CheckAuth, downloadTg *usecase.DownloadTelegramFile, upload *usecase.UploadFile, uploadURL *usecase.UploadURL) *Handler {
	return &Handler{
		api:        api,
		auth:       auth,
		downloadTg: downloadTg,
		upload:     upload,
		uploadURL:  uploadURL,
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
		return
	}

	// check for url in message
	if message.Text != "" {
		h.handleTextMessage(message)
	}
}

func (h *Handler) handleCommand(message *tgbotapi.Message) {
	if message.Command() == "start" {
		h.api.SendText(message.Chat.ID, "any uploaded file will be stored in r2\n\nor you can send direct link for file to upload")
	}
}

func (h *Handler) handleTextMessage(message *tgbotapi.Message) {
	urlRegex := regexp.MustCompile(`https?://[^\s]+`)
	urls := urlRegex.FindAllString(message.Text, -1)

	if len(urls) == 0 {
		return
	}

	// process the first url found
	url := urls[0]
	h.handleURL(message, url)
}

func (h *Handler) handleURL(message *tgbotapi.Message, url string) {
	text, downloadURL, err := h.uploadURL.Execute(context.Background(), url)
	if err != nil {
		h.api.SendText(message.Chat.ID, "failed to upload file from URL: "+err.Error())
		return
	}

	h.api.SendWithButton(message.Chat.ID, text, "✈️", downloadURL)

	if err := h.api.DeleteMessage(message.Chat.ID, message.MessageID); err != nil {
		// intentionally ignored - message deletion is optional cleanup
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
