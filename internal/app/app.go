package app

import (
	"log"

	"github.com/unicodick/r2bot/internal/bot"
	"github.com/unicodick/r2bot/internal/bot/telegram"
	"github.com/unicodick/r2bot/internal/config"
	"github.com/unicodick/r2bot/internal/services"
	"github.com/unicodick/r2bot/internal/storage/r2"
	"github.com/unicodick/r2bot/internal/usecase"
)

type App struct {
	service *bot.Service
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	r2Client, err := r2.NewClient(r2.Config{
		AccountID: cfg.R2AccountID,
		AccessKey: cfg.R2AccessKey,
		SecretKey: cfg.R2SecretKey,
		Bucket:    cfg.R2Bucket,
		PublicURL: cfg.R2PublicURL,
	})
	if err != nil {
		return nil, err
	}

	api, err := telegram.NewAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	downloader := services.NewFileDownloader()

	auth := usecase.NewCheckAuth(cfg)
	downloadTg := usecase.NewDownloadTelegramFile(api.Client(), downloader)
	upload := usecase.NewUploadFile(r2Client)
	uploadURL := usecase.NewUploadURL(r2Client, downloader)

	service := bot.NewService(api, auth, downloadTg, upload, uploadURL)

	return &App{service: service}, nil
}

func (a *App) Run() {
	log.Println("starting r2bot")
	a.service.Start()
}
