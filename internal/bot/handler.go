package bot

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)
const (
	CommandStart         = "/start"
	MessageStart         = "👋 Привет! Чтобы скачать аудио из видео на YouTube, отправь ссылку на видео."
	MessageDownloading   = "⏳ Скачиваю аудио..."
	MessageDownloadError = "❌ Не удалось скачать аудио. Проверьте ссылку и попробуйте снова."
)
func (b *Bot) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	if msg.Text == CommandStart {
		b.sendMessage(chatID, MessageStart)
		return
	}
	b.sendMessage(chatID, MessageDownloading)
	slog.Info("Starting audio download", "url", msg.Text)
	audioPath, title, err := b.downloader.DownloadAudio(ctx, msg.Text)
	if err != nil {
		b.sendMessage(chatID, MessageDownloadError)
		slog.Error("Error download audio", "url", msg.Text, "chatId", chatID, "message", err)
	}
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
