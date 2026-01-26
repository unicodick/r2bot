package telegram

import (
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
