package handler

import (
	"JobPortal/models"
	"JobPortal/service"
	"JobPortal/utils"
	"net/http"
	"strings"

	"github.com/google/uuid"
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

func (h *JobTypeHandlerImpl) GetJobTypeById(c echo.Context) error {
	ctx := c.Request().Context()
	reqID := getRequestID(c)

	id := c.Param("id")

	if id == "" {
		c.Logger().Warnf("GetJobTypeById missing id (request_id=%s)", reqID)
		return utils.BadRequestResponse(c, "id is required", nil)
	}
	jobType, err := h.jobTypeService.GetJobTypeById(ctx, id)

	if err != nil {
		c.Logger().Errorf("GetJobTypeById failed: %v (request_id=%s)", err, reqID)
		return utils.InternalServerErrorResponse(c, "failed to retrieve job type", err.Error())
	}

	if jobType == nil {
		c.Logger().Infof("GetJobTypeById not found: %s (request_id=%s)", id, reqID)
		return utils.NotFoundResponse(c, "job type not found", nil)
	}

	return utils.SuccessResponse(c, http.StatusOK, "Data retrieved successfully", jobType)
}

func (h *JobTypeHandlerImpl) StoreJobType(c echo.Context) (*models.JobType, error) {
	ctx := c.Request().Context()
	reqID := getRequestID(c)

	var input models.CreateJobTypeInput
	if err := c.Bind(&input); err != nil {
		c.Logger().Errorf("StoreJobType bind failed: %v (request_id=%s)", err, reqID)
		return nil, utils.BadRequestResponse(c, "invalid payload", nil)
	}

	if verrs := utils.ValidateStruct(input); verrs != nil {
		c.Logger().Warnf("StoreJobType validation failed (request_id=%s) errors=%v", reqID, verrs)
		return nil, utils.BadRequestResponse(c, "validation failed", verrs)
	}

	newJobType := &models.JobType{
		Name: input.Name,
		Slug: input.Slug,
	}

	created, err := h.jobTypeService.StoreJobType(ctx, newJobType)
	if err != nil {
		c.Logger().Errorf("StoreJobType service error: %v (request_id=%s)", err, reqID)
		return nil, utils.InternalServerErrorResponse(c, "Failed to store job type", err.Error())
	}

	c.Logger().Infof("StoreJobType created id=%s (request_id=%s)", created.Id, reqID)

	_ = utils.SuccessResponse(c, http.StatusCreated, "Job Type Successfully created", created)

	return created, nil
}

func (h *JobTypeHandlerImpl) UpdateJobType(c echo.Context) (*models.JobType, error) {
	ctx := c.Request().Context()
	reqID := getRequestID(c)
	id := strings.TrimSpace(c.Param("id"))

	if id == "" {
		c.Logger().Warnf("UpdateJobType missing id (request_id=%s)", reqID)
		return nil, utils.BadRequestResponse(c, "id is required", nil)
	}

	if _, err := uuid.Parse(id); err != nil {
		c.Logger().Warnf("UpdateJobType invalid id format: %v (request_id=%s)", err, reqID)
		return nil, utils.BadRequestResponse(c, "id must be a valid uuid", nil)
	}

	var input models.CreateJobTypeInput
	if err := c.Bind(&input); err != nil {
		c.Logger().Errorf("UpdateJobType bind failed: %v (request_id=%s)", err, reqID)
		return nil, utils.BadRequestResponse(c, "invalid payload", nil)
	}

	if verrs := utils.ValidateStruct(input); verrs != nil {
		c.Logger().Warnf("UpdateJobType validation failed (request_id=%s) errors=%v", reqID, verrs)
		return nil, utils.BadRequestResponse(c, "validation failed", verrs)
	}

	newJobType := &models.JobType{
		Id:   id,
		Name: input.Name,
		Slug: input.Slug,
	}

	updated, err := h.jobTypeService.UpdateJobType(ctx, newJobType)
	if err != nil {
		c.Logger().Errorf("UpdateJobType service error: %v (request_id=%s)", err, reqID)
		return nil, utils.InternalServerErrorResponse(c, "Failed to update job type", err.Error())
	}

	if updated == nil {
		c.Logger().Infof("UpdateJobType not found id=%s (request_id=%s)", id, reqID)
		return nil, utils.NotFoundResponse(c, "job type not found", nil)
	}

	c.Logger().Infof("UpdateJobType updated id=%s (request_id=%s)", updated.Id, reqID)
	_ = utils.SuccessResponse(c, http.StatusOK, "Job Type Successfully updated", updated)

	return updated, nil
}

func (h *JobTypeHandlerImpl) DeleteJobType(c echo.Context) error {
	ctx := c.Request().Context()
	reqID := getRequestID(c)

	id := strings.TrimSpace(c.Param("id"))
	if _, err := uuid.Parse(id); err != nil {
		c.Logger().Warnf("DeleteJobType invalid id format: %v (request_id=%s)", err, reqID)
		return utils.BadRequestResponse(c, "id must be a valid uuid", nil)
	}

	jobType, err := h.jobTypeService.GetJobTypeById(ctx, id)
	if err != nil {
		c.Logger().Errorf("GetJobType failed: %v (request_id=%s)", err, reqID)
		return utils.NotFoundResponse(c, "job type not found", nil)
	}
	if jobType == nil {
		c.Logger().Infof("GetJobType not found: %s (request_id=%s)", id, reqID)
		return utils.NotFoundResponse(c, "job type not found", nil)
	}

	if err := h.jobTypeService.DeleteJobType(ctx, jobType.Id); err != nil {
		c.Logger().Errorf("DeleteJobType service error: %v (request_id=%s)", err, reqID)
		return utils.InternalServerErrorResponse(c, "Failed to delete job type", err.Error())
	}

	c.Logger().Infof("DeleteJobType id=%s (request_id=%s)", id, reqID)
	return utils.SuccessResponse(c, http.StatusOK, "Job Type Successfully deleted", nil)
}
