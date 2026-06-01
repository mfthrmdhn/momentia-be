package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, envelope{Success: true, Message: message, Data: data})
}

func OKWithMeta(c *gin.Context, message string, data interface{}, meta interface{}) {
	c.JSON(http.StatusOK, envelope{Success: true, Message: message, Data: data, Meta: meta})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, envelope{Success: true, Message: message, Data: data})
}

func BadRequest(c *gin.Context, message string, err interface{}) {
	c.JSON(http.StatusBadRequest, envelope{Success: false, Message: message, Error: err})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, envelope{Success: false, Message: message})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, envelope{Success: false, Message: message})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, envelope{Success: false, Message: message})
}

func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, envelope{Success: false, Message: message})
}

func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, envelope{Success: false, Message: message})
}
