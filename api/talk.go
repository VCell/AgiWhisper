package api

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/VCell/AgiWhisper/module/ctx"
	"github.com/VCell/AgiWhisper/service/agi"
	"github.com/VCell/AgiWhisper/service/transcription"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

type TalkController struct {
	upgrader websocket.Upgrader
}

type AudioFrame struct {
	Action string `json:"action"`
	Audio  string `json:"audio"`
}

type ResponseFrame struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

const (
	ACTION_SPLIT = "split"
	ACTION_ASK   = "ask"

	TYPE_QUESTION = "question"
	TYPE_SLICE    = "slice"
	TYPE_ANSWER   = "answer"
)

func (t *TalkController) TalkManual(c echo.Context) error {
	ctx := ctx.GetCtxFromEcho(c)
	slog.InfoContext(ctx, "TalkManual start")
	ws, err := t.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		slog.ErrorContext(ctx, "Upgrade err", "err", err)
		return err
	}
	defer ws.Close()

	msglist, err := readAudio(ctx, ws)
	if err != nil {
		slog.ErrorContext(ctx, "read audio err:", err)
		return err
	}
	agi := agi.NewAgiSession()
	question, err := agi.ExtractQuestion(ctx, msglist)
	if err != nil {
		slog.ErrorContext(ctx, "extract question err:", err)
		return err
	}
	slog.InfoContext(ctx, "ExtractQuestion:"+question)
	ws.WriteJSON(ResponseFrame{
		Type: TYPE_QUESTION,
		Data: question,
	})
	sliceChannel := make(chan string, 10)
	go agi.AnswerQuestion(ctx, question, sliceChannel)
	var result strings.Builder
	for {
		if slice, ok := <-sliceChannel; ok {
			ws.WriteJSON(ResponseFrame{
				Type: TYPE_SLICE,
				Data: slice,
			})
			result.WriteString(slice)
		} else {
			break
		}
	}
	ws.WriteJSON(ResponseFrame{
		Type: TYPE_ANSWER,
		Data: result.String(),
	})

	return nil
}

func readAudio(ctx context.Context, ws *websocket.Conn) ([]string, error) {
	result := []string{}
	sliceId := 0
	var trans *transcription.AudioTranscription

	defer func() {
		if trans != nil {
			trans.Clean()
		}
	}()
	for {
		if trans == nil {
			trans = transcription.NewAudioTranscription(ctx, sliceId)
			sliceId++
		}
		var frame AudioFrame
		err := ws.ReadJSON(&frame)
		if err != nil {
			slog.ErrorContext(ctx, "error reading frame", "err", err)
			return nil, err
		}
		slog.InfoContext(ctx, fmt.Sprintf("ReadJSON.action:%s,Len:%d", frame.Action, len(frame.Audio)))
		audio, err := base64.StdEncoding.DecodeString(frame.Audio)
		if err != nil {
			slog.ErrorContext(ctx, "error decoding audio", "err", err)
			return nil, err
		}
		if len(audio) > 0 {
			if err = trans.RecordSlice(ctx, audio); err != nil {
				slog.ErrorContext(ctx, "record audio:", "err", err)
				return nil, err
			}
		}

		if frame.Action == ACTION_ASK || frame.Action == ACTION_SPLIT {
			msg, _ := trans.Transcription(ctx)
			trans.Clean()
			trans = nil
			if len(msg) > 0 {
				result = append(result, msg)
			}
			if frame.Action == ACTION_ASK {
				if len(result) > 0 {
					return result, nil
				} else {
					return nil, errors.New("empty question")
				}
			}
		}
	}
}
