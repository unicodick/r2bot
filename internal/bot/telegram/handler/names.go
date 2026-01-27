package handler

import (
	"fmt"
	"regexp"
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
				lowerName := strings.ToLower(archiveName)
				if !strings.HasSuffix(lowerName, ".zip") {
					archiveName += ".zip"
				}
				break
			}
		}
	}

	return archiveName
}

var invalidArchiveNameRe = regexp.MustCompile(`[^a-zA-Z0-9_.-]`)

func (h *Handler) sanitizeArchiveName(name string) string {
	return invalidArchiveNameRe.ReplaceAllString(name, "_")
}
