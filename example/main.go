package main

import (
	"net/http"
	"os"

	"github.com/fitra-besari/pitlog"
	"github.com/labstack/echo/v4"
)

func main() {
	// Buat instance Echo
	e := echo.New()

	// Setup pitlog
	appName := "ExampleApp"
	appVersion := "1.0.0"
	appLevel := "development"
	logDir := "./logs"
	enableLogConsole := "true"
	useSeparate := "false"
	objectView := "false"

	logger, err := pitlog.New_pitlog(appName, appVersion, appLevel, logDir, enableLogConsole, useSeparate, objectView)
	if err != nil {
		e.Logger.Fatal("Failed to initialize logging: ", err)
	}

	// Middleware logging
	logger.Api_log_middleware(e, []string{"Authorization", "password"})

	// Endpoint /ping
	e.GET("/ping", func(c echo.Context) error {
		response := map[string]interface{}{
			"message": "ping!!!",
			"month":   "Server Actived!!",
			"status":  "online",
		}
		return c.JSON(http.StatusOK, response)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
