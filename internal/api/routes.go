package api

import (
	"github.com/labstack/echo/v4"
	"github.com/uchupx/saceri-chatbot-api/internal"
	"github.com/uchupx/saceri-chatbot-api/internal/config"
)

func InitRouter(conf *config.Config, e *echo.Echo) {
	factory := internal.Factory{}

	factory.GetDBConnection()

	middleware := factory.GetMiddleware()

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "Healthy"})
	})

	authRoute := e.Group("/", middleware.Authorization)

	userHandler := factory.GetUserHandler()
	authHandler := factory.GetAuthHandler()

	e.POST("auth/register", authHandler.Register)
	e.POST("auth", authHandler.Login)

	authRoute.GET("user", userHandler.GetUser)
	authRoute.GET("users", userHandler.GetUsers)

	// settings endpoint
	authRoute.GET("settings", factory.GetSettingHandler().GetSetting)
	authRoute.PUT("settings", factory.GetSettingHandler().UpdateSetting)
}
