package bot

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	userID := msg.From.ID

	slog.Info("Processing message",
		"chat_id", chatID,
		"user_id", userID,
		"text", msg.Text,
	)

	slog.Info("Starting audio download", "url", msg.Text)
	b.sendMessage(chatID, "⏳ Скачиваю аудио...")
	slog.Info("Starting audio download", "url", msg.Text)
	audioPath, title, _ := b.downloader.DownloadAudio(ctx, msg.Text)

	audio := tgbotapi.NewAudio(chatID, tgbotapi.FilePath(audioPath))
	audio.Title = title
	audio.Caption = fmt.Sprintf("🎵 %s", title)

	if _, err := b.api.Send(audio); err != nil {
		slog.Error("Failed to send audio", "error", err)
	}
	slog.Info("Finish audio download", "url", msg.Text)
	go DeletePath(audioPath)
}

func (b *Bot) sendMessage(chatID int64, text string) *tgbotapi.Message {
	msg := tgbotapi.NewMessage(chatID, text)

	slog.Debug("Sending message", "text", text)

	sentMsg, err := b.api.Send(msg)
	if err != nil {
		slog.Error("Failed to send message", "error", err)
		return nil
	}

	return &sentMsg
}

func DeletePath(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		slog.Error("Failed to get absolute path", "path", path, "error", err)
		return
	}
	if err := os.Remove(absPath); err != nil {
		slog.Error("Failed to remove audio file", "path", absPath, "error", err)
		return
	}

	slog.Info("Audio file removed", "path", absPath)
}
