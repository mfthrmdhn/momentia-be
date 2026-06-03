package handler

import (
	"crypto/sha256"
	"fmt"
	"momentia-be/pkg/response"
	"momentia-be/repository"
	"momentia-be/services"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo    repository.UserRepository
	service services.UserService
}

func NewUserHandler(repo repository.UserRepository, service services.UserService) *UserHandler {
	return &UserHandler{repo: repo, service: service}
}

func (h *UserHandler) Register(c *gin.Context) {
	var input services.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}

	user, err := h.service.Register(input)
	if err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}
	response.OK(c, "user registered successfully", user)

}

func (h *UserHandler) Login(c *gin.Context) {
	var req services.LoginInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}

	token, err := h.service.Login(req)
	if err != nil {
		response.Unauthorized(c, "invalid credentials")
		return
	}
	response.OK(c, "login successful", gin.H{"token": token})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	user, err := h.repo.GetUserByID(userID.(int))
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}
	response.OK(c, "user retrieved successfully", user)
}

func (h *UserHandler) Logout(c *gin.Context) {
	tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if tokenStr == "" {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	hash := sha256.Sum256([]byte(tokenStr))
	tokenHash := fmt.Sprintf("%x", hash)

	if err := h.service.Logout(tokenHash); err != nil {
		response.InternalError(c, "logout failed")
		return
	}
	response.OK(c, "logout successful", nil)
}