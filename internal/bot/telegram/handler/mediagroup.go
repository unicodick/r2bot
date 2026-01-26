package telegram

import (
	"context"

	"github.com/unicodick/r2bot/internal/archiver"
	telegramgroup "github.com/unicodick/r2bot/internal/bot/telegram/mediagroup"
)

func (h *Handler) handleMediaGroup(group *telegramgroup.MediaGroup) {
	if len(group.Files) == 0 {
		return
	}

	chatID := group.ChatID

	// if single file, handle as normal upload
	if len(group.Files) == 1 {
		h.handleSingleFile(group.Files[0], chatID)
		return
	}

	// multiple files - create archive
	h.handleMultipleFiles(group)
}

func (h *Handler) handleSingleFile(file telegramgroup.MediaGroupFile, chatID int64) {
	text, url, err := h.upload.Execute(context.Background(), file.FileData.Name, file.FileData.Reader, file.FileData.Size)
	if err != nil {
		h.api.SendText(chatID, "failed to upload file")
		return
	}

	h.api.SendWithButton(chatID, text, "✈️", url)

	if err := h.api.DeleteMessage(chatID, file.Message.MessageID); err != nil {
		// intentionally ignored - message deletion is optional cleanup
	}
}

func (h *Handler) handleMultipleFiles(group *telegramgroup.MediaGroup) {
	chatID := group.ChatID

	// prepare files for archiving
	archiveFiles := make([]archiver.FileData, len(group.Files))
	for i, file := range group.Files {
		archiveFiles[i] = file.FileData
	}

	archiveName := h.generateArchiveName(group)

	text, url, err := h.uploadArchive.Execute(context.Background(), archiveFiles, archiveName)
	if err != nil {
		h.api.SendText(chatID, "failed to create and upload archive: "+err.Error())
		return
	}

	h.api.SendWithButton(chatID, text, "✈️", url)

	h.deleteGroupMessages(group)
}

func (h *Handler) deleteGroupMessages(group *telegramgroup.MediaGroup) {
	for _, file := range group.Files {
		if err := h.api.DeleteMessage(group.ChatID, file.Message.MessageID); err != nil {
			// intentionally ignored - message deletion is optional cleanup
		}
	}
}
