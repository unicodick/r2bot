package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// api is thin wrapper over telegram sdk
type API struct {
	client *Client
	sender *Sender
}

func NewAPI(token string) (*API, error) {
	client, err := NewClient(token)
	if err != nil {
		return nil, err
	}

	sender := NewSender(client)

	return &API{
		client: client,
		sender: sender,
	}, nil
}

func (a *API) GetUpdates() tgbotapi.UpdatesChannel {
	return a.client.GetUpdates()
}

func (a *API) SendText(chatID int64, text string) error {
	return a.sender.SendText(chatID, text)
}

func (a *API) SendWithButton(chatID int64, text, buttonText, url string) error {
	return a.sender.SendWithButton(chatID, text, buttonText, url)
}

func (a *API) DeleteMessage(chatID int64, messageID int) error {
	return a.client.Send(tgbotapi.NewDeleteMessage(chatID, messageID))
}

func (a *API) GetUsername() string {
	return a.client.GetUsername()
}

// client returns underlying client for download operations
func (a *API) Client() *Client {
	return a.client
}
