package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/unicodick/r2bot/internal/config"
	"github.com/unicodick/r2bot/internal/storage"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	config *config.Config
	r2     *storage.R2Client
}

func New(cfg *config.Config, r2Client *storage.R2Client) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:    api,
		config: cfg,
		r2:     r2Client,
	}, nil
}

func (b *Bot) Start() {
	log.Printf("bot started as @%s", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		if !b.config.IsUserAllowed(userID) {
			continue
		}

		b.handleMessage(update.Message)
	}
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	b.api.Send(msg)
}

func (b *Bot) sendMessageWithButton(chatID int64, text, url string) {
	msg := tgbotapi.NewMessage(chatID, text)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonSwitch("✈️", url),
		),
	)

	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}
