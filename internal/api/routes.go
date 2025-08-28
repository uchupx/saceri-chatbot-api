package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uchupx/saceri-chatbot-api/internal"
	"github.com/uchupx/saceri-chatbot-api/internal/config"
)

func InitRouter(conf *config.Config, e *echo.Echo) {
	factory := internal.Factory{}

	factory.GetDBConnection()

	apiMiddleware := factory.GetMiddleware()

	optionsHandler := func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		return c.NoContent(http.StatusNoContent)
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodOptions},
	}))

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "Healthy"})
	})
	e.OPTIONS("/health", optionsHandler)

	authRoute := e.Group("/", apiMiddleware.Authorization)

	userHandler := factory.GetUserHandler()
	authHandler := factory.GetAuthHandler()

	e.POST("auth/register", authHandler.Register)
	e.OPTIONS("auth/register", optionsHandler)

	e.POST("auth", authHandler.Login)
	e.OPTIONS("auth", optionsHandler)

	authRoute.GET("user", userHandler.GetUser)
	authRoute.OPTIONS("user", optionsHandler)

	authRoute.GET("user/me", userHandler.GetMe)
	authRoute.PUT("user/me", userHandler.UpdateMe)
	authRoute.OPTIONS("user/me", optionsHandler)

	authRoute.GET("users", userHandler.GetUsers)
	authRoute.OPTIONS("users", optionsHandler)

	authRoute.GET("settings", factory.GetSettingHandler().GetSetting)
	authRoute.PUT("settings", factory.GetSettingHandler().UpdateSetting)
	authRoute.OPTIONS("settings", optionsHandler)
}
