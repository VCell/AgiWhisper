package agi

import (
	"context"
	"log/slog"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type AgiAnswer struct {
	Key  string
	Info string
}

type AgiSession struct {
}

func NewAgiSession() *AgiSession {
	return &AgiSession{}
}

func (a *AgiSession) ExtractQuestion(ctx context.Context, msglist []string) (string, error) {
	llm, err := openai.New(openai.WithModel(os.Getenv("AGI_MODEL")))
	if err != nil {
		slog.ErrorContext(ctx, "new gpt error:", err)
		return "", err
	}
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, PROMPT_SYSTEM_EXTRACT_QUESTION),
		llms.TextParts(llms.ChatMessageTypeHuman, msglist...),
	}
	completion, err := llm.GenerateContent(ctx, content)
	if err != nil {
		slog.ErrorContext(ctx, "gpt generate error:", err)
		return "", err
	}
	return completion.Choices[0].Content, nil
}

func (a *AgiSession) AnswerQuestion(ctx context.Context, question string, sliceChan chan string) error {
	defer close(sliceChan)
	llm, err := openai.New(openai.WithModel(os.Getenv("AGI_MODEL")))
	if err != nil {
		slog.ErrorContext(ctx, "new gpt error:"+err.Error())
		return err
	}
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, PROMPT_SYSTEM_ANSWER_QUESTION),
		llms.TextParts(llms.ChatMessageTypeHuman, question),
	}
	completion, err := llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		sliceChan <- string(chunk)
		return nil
	}))
	if err != nil {
		slog.ErrorContext(ctx, "gpt generate error:", err)
		return err
	}
	slog.InfoContext(ctx, "GenerateContent:"+completion.Choices[0].Content)
	return nil
}
