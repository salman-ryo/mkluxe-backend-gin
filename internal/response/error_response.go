package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Error(c *gin.Context, status int, errorCode, message string, fieldErrors []FieldError) {
	payload := gin.H{
		"success": false,
		"code":    errorCode,
		"message": message,
	}
	if len(fieldErrors) > 0 {
		payload["errors"] = fieldErrors
	}
	c.JSON(status, payload)
}

func BadRequest(c *gin.Context, message string, errors []FieldError) {
	Error(c, http.StatusBadRequest, "BAD_REQUEST", message, errors)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, "FORBIDDEN", message, nil)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, "NOT_FOUND", message, nil)
}

func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, nil)
}
