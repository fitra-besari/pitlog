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
	enableLogConsole := "1" // "0"  for enable log on conslosso ,default  will be outo enable on log file only
	useSeparate := "1"      // "0" if use 1 you can ssee separate border on every request
	objectView := "false"   // "0" if use 1 you can see log data on object json if use 0 you can see on string data

	//  Create instance
	logger, err := pitlog.New_pitlog(appName, appVersion, appLevel, logDir, enableLogConsole, useSeparate, objectView)
	if err != nil {
		e.Logger.Fatal("Failed to initialize logging: ", err)
	}

	// Register Middleware logging
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
