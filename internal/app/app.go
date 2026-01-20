package app

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/unicodick/r2bot/internal/bot"
	"github.com/unicodick/r2bot/internal/config"
	"github.com/unicodick/r2bot/internal/storage"
)

type App struct {
	bot *bot.Bot
}

func New() (*App, error) {
	godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	r2Client, err := storage.NewR2Client(
		cfg.R2AccountID,
		cfg.R2AccessKey,
		cfg.R2SecretKey,
		cfg.R2Bucket,
		cfg.R2PublicURL,
	)
	if err != nil {
		return nil, err
	}

	bot, err := bot.New(cfg, r2Client)
	if err != nil {
		return nil, err
	}

	return &App{bot: bot}, nil
}

func (a *App) Run() {
	log.Println("starting r2bot")
	a.bot.Start()
}
