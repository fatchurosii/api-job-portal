package handler

import (
	"JobPortal/service"
	"JobPortal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type JobTypeHandlerImpl struct {
	jobTypeService service.JobTypeService
}

func NewJobTypeHandler(jobTypeService service.JobTypeService) *JobTypeHandlerImpl {
	return &JobTypeHandlerImpl{jobTypeService: jobTypeService}
}

func getRequestID(c echo.Context) string {
	if v := c.Response().Header().Get(echo.HeaderXRequestID); v != "" {
		return v
	}
	if v := c.Request().Header.Get(echo.HeaderXRequestID); v != "" {
		return v
	}
	if v := c.Get("request_id"); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (h *JobTypeHandlerImpl) GetAllJobType(c echo.Context) error {

	ctx := c.Request().Context()

	jobTypes, err := h.jobTypeService.GetAllJobTypes(ctx)
	reqID := getRequestID(c)

	if err != nil {
		c.Logger().Errorf("GetAllJobType failed: %v (request_id=%s)", err, reqID)
		return utils.InternalServerErrorResponse(c, "Failed to retrieve job types", err.Error())
	}

	// log success
	c.Logger().Infof("GetAllJobType success: returned %d items (request_id=%s)", len(jobTypes), reqID)

	return utils.SuccessResponse(c, http.StatusOK, "Data retrieve successfully", jobTypes)

}
