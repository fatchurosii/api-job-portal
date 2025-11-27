package utils

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Response struct {
	RequestId string      `json:"requestId,omitempty"`
	Status    string      `json:"status"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
}

func SuccessResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	reqID := uuid.New().String()
	response := Response{
		RequestId: reqID,
		Status:    "success",
		Message:   message,
		Data:      data,
	}

	c.Response().Header().Set(echo.HeaderXRequestID, reqID)

	return c.JSON(statusCode, response)
}

func errorResponse(c echo.Context, statusCode int, message string, errors interface{}) error {
	reqID := uuid.New().String()
	response := Response{
		RequestId: reqID,
		Status:    "error",
		Message:   message,
		Data:      nil,
		Errors:    errors,
	}

	c.Response().Header().Set(echo.HeaderXRequestID, reqID)

	return c.JSON(statusCode, response)
}

func badRequestResponse(c echo.Context, message string, errors interface{}) error {
	return errorResponse(c, http.StatusBadRequest, message, errors)
}

func unauthorizedResponse(c echo.Context, message string, errors interface{}) error {
	return errorResponse(c, http.StatusUnauthorized, message, errors)
}

func notFoundResponse(c echo.Context, message string, errors interface{}) error {
	return errorResponse(c, http.StatusNotFound, message, errors)
}

func forbiddenResponse(c echo.Context, message string, errors interface{}) error {
	return errorResponse(c, http.StatusForbidden, message, errors)
}

func InternalServerErrorResponse(c echo.Context, message string, err interface{}) error {
	return errorResponse(c, http.StatusInternalServerError, message, err)
}

func unprocessableEntityResponse(c echo.Context, message string, errors interface{}) error {
	return errorResponse(c, http.StatusUnprocessableEntity, message, errors)
}
