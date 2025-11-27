package routes

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

func RegisterAllRoutes(e *echo.Echo, db *sql.DB) {
	api := e.Group("/api")

	v1 := api.Group("/v1")

	registerJobTypeRoute(v1, db)
}
