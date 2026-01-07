package downloader

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"soundExtractBot/internal/utils"

	"github.com/kkdai/youtube/v2"
)

const DownloadDir = "./downloads"

type YouTubeDownloader struct {
	client *youtube.Client
}

func NewYouTubeDownloader() *YouTubeDownloader {
	return &YouTubeDownloader{
		client: &youtube.Client{},
	}
}

func (d *YouTubeDownloader) DownloadAudio(
	ctx context.Context,
	url string,
) (filePath string, title string, err error) {

	slog.Info("Downloading YouTube audio", "url", url)

	video, err := d.client.GetVideoContext(ctx, url)
	if err != nil {
		return "", "", fmt.Errorf("get video: %w", err)
	}

	formats := video.Formats.WithAudioChannels()
	if len(formats) == 0 {
		return "", "", fmt.Errorf("no audio formats available")
	}

	audioFormat := utils.GetAudioFormat(formats)
	stream, _, err := d.client.GetStreamContext(ctx, video, &audioFormat)
	if err != nil {
		return "", "", fmt.Errorf("get stream: %w", err)
	}
	defer func() {
		if err := stream.Close(); err != nil {
			slog.Error("Failed to close stream", "error", err)
		}
	}()

	ext := "m4a"
	if audioFormat.MimeType == "audio/webm" {
		ext = "webm"
	}

	// создаём директорию, если её нет
	if err := os.MkdirAll(DownloadDir, os.ModePerm); err != nil {
		return "", "", fmt.Errorf("create download dir: %w", err)
	}

	filePath = filepath.Join(DownloadDir, fmt.Sprintf("%s.%s", video.ID, ext))

	file, err := os.Create(filePath)
	if err != nil {
		return "", "", fmt.Errorf("create file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("Failed to close file", "error", err)
		}
	}()

	if _, err := io.Copy(file, stream); err != nil {
		return "", "", fmt.Errorf("save audio: %w", err)
	}

	slog.Info("Audio downloaded",
		"title", video.Title,
		"file", filePath,
		"mime", audioFormat.MimeType,
		"bitrate", audioFormat.Bitrate,
	)

	return filePath, video.Title, nil
}
