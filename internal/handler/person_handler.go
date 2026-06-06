package handler

import (
	"errors"
	"momentia-be/pkg/response"
	"momentia-be/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PersonHandler struct {
	service services.PersonService
}

func NewPersonHandler(service services.PersonService) *PersonHandler {
	return &PersonHandler{service: service}
}

func (h *PersonHandler) CreatePerson(c *gin.Context) {
	var input services.CreatePersonInput
	userID := c.GetInt("userID")
	if userID == 0 {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}

	person, err := h.service.CreatePerson(userID, input)
	if err != nil {
		response.InternalError(c, "failed to create person")
		return
	}

	response.OK(c, "person created successfully", person)
}

func (h *PersonHandler) GetPersonByID(c *gin.Context) {
	parsedID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid person ID", nil)
		return
	}

	person, err := h.service.GetPersonByID(parsedID, c.GetInt("userID"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "person not found")
			return
		}
		response.InternalError(c, "failed to retrieve person")
		return
	}
	response.OK(c, "person retrieved successfully", person)
}

func (h *PersonHandler) GetPersons(c *gin.Context) {
	userID := c.GetInt("userID")
	if userID == 0 {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	persons, meta, err := h.service.GetAllPersons(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, "failed to retrieve persons")
		return
	}
	response.OK(c, "persons retrieved successfully", gin.H{
		"data":       persons,
		"pagination": meta,
	})
}

func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	var input services.UpdatePersonInput
	userID := c.GetInt("userID")
	if userID == 0 {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	parsedID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid person ID", nil)
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}

	person, err := h.service.UpdatePerson(parsedID, userID, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "person not found")
			return
		}
		response.InternalError(c, "failed to update person")
		return
	}
	response.OK(c, "person updated successfully", person)
}

func (h *PersonHandler) DeletePerson(c *gin.Context) {
	parsedID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid person ID", nil)
		return
	}

	if err := h.service.DeletePerson(parsedID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "person not found")
			return
		}
		response.InternalError(c, "failed to delete person")
		return
	}
	response.OK(c, "person deleted successfully", nil)
}
