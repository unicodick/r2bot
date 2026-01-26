package handler

import (
	"fmt"
	"strings"
	"time"

	telegramgroup "github.com/unicodick/r2bot/internal/bot/telegram/mediagroup"
)

func (h *Handler) generateArchiveName(group *telegramgroup.MediaGroup) string {
	archiveName := fmt.Sprintf("files_%d.zip", time.Now().Unix())

	for _, file := range group.Files {
		if file.Message.Caption != "" {
			caption := strings.TrimSpace(file.Message.Caption)
			if caption != "" {
				archiveName = h.sanitizeArchiveName(caption)
				if !strings.HasSuffix(strings.ToLower(archiveName), ".zip") {
					archiveName += ".zip"
				}
				break
			}
		}
	}

	return archiveName
}

func (h *Handler) sanitizeArchiveName(name string) string {
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, ":", "_")
	name = strings.ReplaceAll(name, "*", "_")
	name = strings.ReplaceAll(name, "?", "_")
	name = strings.ReplaceAll(name, "\"", "_")
	name = strings.ReplaceAll(name, "<", "_")
	name = strings.ReplaceAll(name, ">", "_")
	name = strings.ReplaceAll(name, "|", "_")

	return name
}
