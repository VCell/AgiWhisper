package main

import (
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"

	"github.com/VCell/AgiWhisper/api"
	"github.com/VCell/AgiWhisper/module/log"
	"github.com/labstack/echo/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Load env error:", err)
		return
	}
	talkCtr := api.TalkController{}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(log.TraceIDMiddleware)
	e.Static("/web", "web")
	e.GET("/talk_manual", talkCtr.TalkManual)
	e.Logger.Fatal(e.Start(":8000"))
}
