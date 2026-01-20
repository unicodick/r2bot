package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// client is thin wrapper over telegram bot api
type Client struct {
	bot *tgbotapi.BotAPI
}

func NewClient(token string) (*Client, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Client{bot: bot}, nil
}

func (c *Client) GetUpdates() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return c.bot.GetUpdatesChan(u)
}

func (c *Client) Send(msg tgbotapi.Chattable) error {
	_, err := c.bot.Send(msg)
	return err
}

func (c *Client) GetUsername() string {
	return c.bot.Self.UserName
}

func (c *Client) GetFileURL(fileID string) (string, error) {
	file, err := c.bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return "", err
	}

	return file.Link(c.bot.Token), nil
}
