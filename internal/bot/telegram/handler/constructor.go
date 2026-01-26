package telegram

import (
	"context"
	"time"

	telegramapi "github.com/unicodick/r2bot/internal/bot/telegram"
	telegramgroup "github.com/unicodick/r2bot/internal/bot/telegram/mediagroup"
	"github.com/unicodick/r2bot/internal/usecase"
)

func NewHandler(api *telegramapi.API, auth *usecase.CheckAuth, downloadTg *usecase.DownloadTelegramFile, upload *usecase.UploadFile, uploadURL *usecase.UploadURL, uploadArchive *usecase.UploadArchive) *Handler {
	h := &Handler{
		api:           api,
		auth:          auth,
		downloadTg:    downloadTg,
		upload:        upload,
		uploadURL:     uploadURL,
		uploadArchive: uploadArchive,
	}

	// init mediagroup
	h.mediaGroupMgr = telegramgroup.NewMediaGroupManager(5*time.Second, h.handleMediaGroup)

	go h.mediaGroupMgr.StartCleanupRoutine(context.Background(), 1*time.Minute, 10*time.Minute)

	return h
}
