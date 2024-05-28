package transcription

import (
	"context"
	"log/slog"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

const AUDIO_EXT = "m4a"

type AudioTranscription struct {
	filepath string
}

func NewAudioTranscription() *AudioTranscription {
	filename := uuid.New().String()
	basepath := os.Getenv("AUDIO_RECORD_PATH")
	return &AudioTranscription{
		filepath: path.Join(basepath, filename+"."+AUDIO_EXT),
	}
}

func (a *AudioTranscription) RecordSlice(ctx context.Context, data []byte) error {
	file, err := os.OpenFile(a.filepath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		slog.ErrorContext(ctx, "file open error:", err)
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		slog.ErrorContext(ctx, "file write error:", err)
		return err
	}
	return nil
}

func (a *AudioTranscription) Transcription(ctx context.Context) (string, error) {
	defer a.Clean()
	config := openai.DefaultConfig(os.Getenv("OPENAI_API_KEY"))
	config.BaseURL = os.Getenv("OPENAI_BASE_URL")

	client := openai.NewClientWithConfig(config)
	resp, err := client.CreateTranscription(
		context.Background(),
		openai.AudioRequest{
			Model:    openai.Whisper1,
			FilePath: a.filepath,
			Prompt:   "使用简体中文输出",
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "transcription error:", err)
		return "", err
	}
	// if len(resp.Text) == 0 {
	// 	slog.ErrorContext(ctx, "transcription result is empty.")
	// 	return "", errors.New("transcription result is empty")
	// }
	slog.InfoContext(ctx, "Transcription result:"+resp.Text)
	return resp.Text, nil
}

func (a *AudioTranscription) Clean() error {
	if os.Getenv("AUTO_DELETE_AUDIO_RECORD") == "true" {
		return os.Remove(a.filepath)
	}
	return nil
}
