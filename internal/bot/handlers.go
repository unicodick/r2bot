package bot

import (
	"fmt"
	"io"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	if message.IsCommand() {
		b.handleCommand(message)
		return
	}

	if message.Document != nil {
		b.handleDocument(message)
		return
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		b.sendMessage(message.Chat.ID, "any uploaded file will be stored in r2")
	}
}

func (b *Bot) handleDocument(message *tgbotapi.Message) {
	document := message.Document

	if document.FileSize > 50*1024*1024 { // 50MB limit
		b.sendMessage(message.Chat.ID, "file too big (<50mb)")
		return
	}

	fileReader, fileSize, err := b.downloadTelegramFile(document.FileID)
	if err != nil {
		log.Printf("failed to download file: %v", err)
		b.sendMessage(message.Chat.ID, "failed to upload (dwnld)")
		return
	}
	defer fileReader.Close()

	filename := document.FileName
	if filename == "" {
		filename = "file"
	}

	fileInfo, err := b.r2.UploadFile(filename, fileReader, fileSize)
	if err != nil {
		log.Printf("failed to upload to r2: %v", err)
		b.sendMessage(message.Chat.ID, "failed to upload (r2)")
		return
	}

	text := fmt.Sprintf("file: %s\nsize: %s\nlink: %s",
		fileInfo.Name,
		formatFileSize(fileInfo.Size),
		fileInfo.URL)

	b.sendMessageWithButton(message.Chat.ID, text, fileInfo.URL)
}

func (b *Bot) downloadTelegramFile(fileID string) (io.ReadCloser, int64, error) {
	fileConfig := tgbotapi.FileConfig{FileID: fileID}
	file, err := b.api.GetFile(fileConfig)
	if err != nil {
		return nil, 0, err
	}

	fileURL := file.Link(b.api.Token)
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, 0, err
	}

	return resp.Body, resp.ContentLength, nil
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d b", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cb", float64(bytes)/float64(div), "kmgtpe"[exp])
}
