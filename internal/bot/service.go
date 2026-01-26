package bot

import (
	"log"

	"github.com/unicodick/r2bot/internal/bot/telegram"
	telegramhandler "github.com/unicodick/r2bot/internal/bot/telegram/handler"
	"github.com/unicodick/r2bot/internal/usecase"
)

type Service struct {
	api     *telegram.API
	handler *telegramhandler.Handler
}

func NewService(api *telegram.API, auth *usecase.CheckAuth, downloadTg *usecase.DownloadTelegramFile, upload *usecase.UploadFile, uploadURL *usecase.UploadURL, uploadArchive *usecase.UploadArchive) *Service {
	handler := telegramhandler.NewHandler(api, auth, downloadTg, upload, uploadURL, uploadArchive)

	return &Service{
		api:     api,
		handler: handler,
	}
}

func (s *Service) Start() {
	log.Printf("bot started as @%s", s.api.GetUsername())

	updates := s.api.GetUpdates()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		s.handler.HandleMessage(update.Message)
	}
}
