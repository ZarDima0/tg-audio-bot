package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"soundExtractBot/internal/bot"
	"strconv"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	logLevel := getLogLevel()
	jsonLogs := getJSONLogs()
	setupLogger(logLevel, jsonLogs)
	slog.Info("Starting YouTube Audio Bot",
		"log_level", logLevel,
		"json_logs", jsonLogs,
	)
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		slog.Error("TELEGRAM_TOKEN not set")
		os.Exit(1)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		slog.Info("Received shutdown signal, stopping bot...")
		cancel()
	}()
	botInstance, err := bot.NewBot(token)
	if err != nil {
		slog.Error("Failed to create bot", "error", err)
		os.Exit(1)
	}

	botInstance.Starting(ctx)
	slog.Info("Bot started successfully")
}

func getLogLevel() slog.Level {
	levelStr := os.Getenv("LOG_LEVEL")
	switch levelStr {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func getJSONLogs() bool {
	jsonLogs, _ := strconv.ParseBool(os.Getenv("JSON_LOGS"))
	return jsonLogs
}

func setupLogger(level slog.Level, jsonLogs bool) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: level,
	}

	if jsonLogs {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}
