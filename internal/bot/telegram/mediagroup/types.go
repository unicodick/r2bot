package telegram

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/unicodick/r2bot/internal/archiver"
)

type MediaGroupFile struct {
	Message  *tgbotapi.Message
	FileData archiver.FileData
}

type MediaGroup struct {
	ID      string
	ChatID  int64
	Files   []MediaGroupFile
	Timer   *time.Timer
	Created time.Time
}
