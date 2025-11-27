package main

import (
	"JobPortal/routes"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"JobPortal/config"
	"JobPortal/database"
	"JobPortal/utils"

	"github.com/labstack/echo/v4"
)

func main() {

	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatal("Error loading environment variables: ", err)
	}

	// Connect to database
	if err := database.InitDb(); err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	defer func() {
		if err := database.DB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return utils.SuccessResponse(c, http.StatusOK, "Hallo World!", nil)
	})

	e.GET("/health-check", func(c echo.Context) error {
		data, statusCode := healthCheck(c.Request().Context())
		return utils.SuccessResponse(c, statusCode, "Health Check", data)
	})

	routes.RegisterAllRoutes(e, database.DB)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server %s on port %s", os.Getenv("APP_NAME"), port)

	e.Logger.Fatal(e.Start(":" + port))
}

func healthCheck(ctx context.Context) (map[string]interface{}, int) {
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "JobPortal"
	}
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	dbStatus := "unknown"
	var dbErrStr string
	var dbLatencyMs int64 = -1

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	start := time.Now()
	if err := database.DB.PingContext(pingCtx); err != nil {
		dbStatus = "unavailable"
		dbErrStr = err.Error()
	} else {
		dbStatus = "ok"
		dbLatencyMs = time.Since(start).Milliseconds()
	}

	// Aggregate response data
	data := map[string]interface{}{
		"app": map[string]interface{}{
			"name": appName,
			"env":  appEnv,
		},
		"config_loaded": true,
		"database": map[string]interface{}{
			"status":  dbStatus,
			"error":   dbErrStr,
			"ping_ms": dbLatencyMs,
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	// Jika DB down, kembalikan 503; selain itu 200
	statusCode := http.StatusOK
	if dbStatus != "ok" {
		statusCode = http.StatusServiceUnavailable
	}

	return data, statusCode
}
