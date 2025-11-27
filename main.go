package main

import (
	"JobPortal/config"
	"JobPortal/database"
	"log"
	"net/http"
	"os"

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
		return c.String(http.StatusOK, "Hello, World!")
	})

	port := os.Getenv("APP_PORT")
	appName := os.Getenv("APP_NAME")

	if appName == "" {
		appName = "JobPortal"
	}

	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server %s on port %s", appName, port)

	e.Logger.Fatal(e.Start(": " + port))
}
