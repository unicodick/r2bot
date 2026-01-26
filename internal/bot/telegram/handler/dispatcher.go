package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
