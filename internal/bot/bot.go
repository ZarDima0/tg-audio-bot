package bot

import (
	"context"
	"log/slog"
	"soundExtractBot/internal/downloader"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api        *tgbotapi.BotAPI
	downloader *downloader.YouTubeDownloader
}

func NewBot(token string) (*Bot, error) {
	slog.Info("Initializing Telegram bot API")

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		slog.Error("Failed to create bot API", "key", err)
	}

	slog.Info("Bot authorized", "username", api.Self.UserName)

	dl := downloader.NewYouTubeDownloader()
	return &Bot{
		api:        api,
		downloader: dl,
	}, nil
}

func (b *Bot) Starting(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	slog.Info("Starting updates polling")

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Shutting down bot")
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			slog.Info("Received message",
				"chat_id", update.Message.Chat.ID,
				"user_id", update.Message.From.ID,
				"text", update.Message.Text,
			)
			go b.handleMessage(ctx, update.Message)
		}
	}
}
