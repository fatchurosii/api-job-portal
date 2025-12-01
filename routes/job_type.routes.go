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
	group.GET("/:id", h.GetJobTypeById)
	group.POST("", func(c echo.Context) error {
		_, err := h.StoreJobType(c)
		return err
	})
	group.PUT("/:id", func(c echo.Context) error {
		_, err := h.UpdateJobType(c)
		return err
	})
	group.DELETE("/:id", h.DeleteJobType)

	group.PATCH("/:id/change-status", func(c echo.Context) error {
		_, err := h.ChangeStatusJobType(c)
		return err
	})

}
