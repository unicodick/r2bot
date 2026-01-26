package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleCommand(message *tgbotapi.Message) {
	if message.Command() == "start" {
		h.api.SendText(message.Chat.ID, "any uploaded file will be stored in r2\n\nor you can send direct link for file to upload")
	}
}
