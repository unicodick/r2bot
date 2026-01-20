package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Sender struct {
	client *Client
}

func NewSender(client *Client) *Sender {
	return &Sender{
		client: client,
	}
}

func (s *Sender) SendText(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	return s.client.Send(msg)
}

func (s *Sender) SendWithButton(chatID int64, text, buttonText, url string) error {
	msg := tgbotapi.NewMessage(chatID, text)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonSwitch(buttonText, url),
		),
	)

	msg.ReplyMarkup = keyboard
	return s.client.Send(msg)
}
