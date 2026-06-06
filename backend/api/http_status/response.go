package http_status

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type successResponse struct {
	Data any `json:"data"`
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type listResponse struct {
	Data       any `json:"data"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

// Success sends a 200 JSON response with the given data.
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, successResponse{Data: data})
}

// Created sends a 201 JSON response with the given data.
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, successResponse{Data: data})
}

// NoContent sends a 204 response with no body.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// List sends a paginated list response.
func List(c *gin.Context, data any, total, page, limit int) {
	totalPages := 0
	if limit > 0 && total > 0 {
		totalPages = (total + limit - 1) / limit
	}
	c.JSON(http.StatusOK, listResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// Error sends an error JSON response, deriving the HTTP status from the error type.
func Error(c *gin.Context, err error) {
	status := HTTPStatusFromError(err)
	code := CodeFromError(err)
	c.JSON(status, errorResponse{Code: code, Message: err.Error()})
}

// BadRequest sends a 400 response.
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, errorResponse{Code: "BAD_REQUEST", Message: msg})
}

// Unauthorized sends a 401 response.
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, errorResponse{Code: "UNAUTHORIZED", Message: msg})
}

// Forbidden sends a 403 response.
func Forbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, errorResponse{Code: "FORBIDDEN", Message: msg})
}

// NotFound sends a 404 response.
func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, errorResponse{Code: "NOT_FOUND", Message: msg})
}

// InternalError sends a 500 response.
func InternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, errorResponse{Code: "INTERNAL_ERROR", Message: "internal server error"})
}
