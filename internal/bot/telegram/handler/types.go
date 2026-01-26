package telegram

import (
	telegramapi "github.com/unicodick/r2bot/internal/bot/telegram"
	telegramgroup "github.com/unicodick/r2bot/internal/bot/telegram/mediagroup"
	"github.com/unicodick/r2bot/internal/usecase"
)

type Handler struct {
	api           *telegramapi.API
	auth          *usecase.CheckAuth
	downloadTg    *usecase.DownloadTelegramFile
	upload        *usecase.UploadFile
	uploadURL     *usecase.UploadURL
	uploadArchive *usecase.UploadArchive
	mediaGroupMgr *telegramgroup.MediaGroupManager
}
