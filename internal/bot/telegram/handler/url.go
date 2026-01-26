package telegram

import (
	"context"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleTextMessage(message *tgbotapi.Message) {
	urlRegex := regexp.MustCompile(`https?://[^\s]+`)
	urls := urlRegex.FindAllString(message.Text, -1)

	if len(urls) == 0 {
		return
	}

	// process the first url found
	url := strings.TrimRight(urls[0], ".,!?;:)'\"")
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
