package routes

import (
	"JobPortal/handler"
	"JobPortal/repository"
	"JobPortal/service"
	"database/sql"

	"github.com/labstack/echo/v4"
)

func registerJobTypeRoute(e *echo.Group, db *sql.DB) {
	repo := repository.NewJobTypeRepository(db)
	svc := service.NewJobTypeService(repo)
	h := handler.NewJobTypeHandler(svc)

	group := e.Group("/job-types")
	group.GET("", h.GetAllJobType)
}
