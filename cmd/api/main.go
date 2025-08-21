package main

import (
	"github.com/labstack/echo/v4"
	"github.com/uchupx/saceri-chatbot-api/internal/api"
	"github.com/uchupx/saceri-chatbot-api/internal/config"
)

func main() {
	e := echo.New()
	e.Debug = true

	conf := config.GetConfig()

	runAPIServer(conf, e)
}

func runAPIServer(conf *config.Config, e *echo.Echo) {
	api.InitRouter(conf, e)
	if err := e.Start(":" + conf.App.Port); err != nil {
		e.Logger.Fatal(err.Error())
	}
}
