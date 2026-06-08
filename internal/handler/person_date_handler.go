package handler

import (
	"errors"
	"momentia-be/pkg/response"
	"momentia-be/services"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type PersonDateHandler struct {
	service services.PersonDatesService
}

func NewPersonDateHandler(service services.PersonDatesService) *PersonDateHandler {
	return &PersonDateHandler{service: service}
}

func (h *PersonDateHandler) CreatePersonDate(c *gin.Context) {
	var input services.CreatePersonDatesInput
	userID := c.GetInt("userID")
	if userID == 0 {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	pid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid person id", err.Error())
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}

	dates, err := h.service.CreatePersonDates(pid, input)
	if err != nil {
		response.InternalError(c, "failed to create dates")
		return
	}

	response.OK(c, "dates created successfully", dates)
}

func (h *PersonDateHandler) GetAllPersonDates(c *gin.Context) {
	userID := c.GetInt("userID")
	if userID == 0 {
		response.Unauthorized(c, "user not authenticated")
		return
	}

	pid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid person id", err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	dates, meta, err := h.service.GetAllPersonDates(pid, page, pageSize)
	if err != nil {
		response.InternalError(c, "failed to retrieve dates")
		return
	}
	response.OK(c, "dates retrieved successfully", gin.H{
		"data": dates,
		"pagination": meta,
	})
}

func (h *PersonDateHandler) GetPersonDatesByID(c *gin.Context) {
	parsedPersonID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid person ID", nil)
		return
	}

	parsedDateID, err := uuid.Parse(c.Param("dateId"))
	if err != nil {
		response.BadRequest(c, "invalid dates ID", nil)
		return
	}
	
	personDates, err := h.service.GetPersonDatesByID(parsedPersonID, parsedDateID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "person / dates not found")
			return
		}
		response.InternalError(c, "failed to retrieve dates")
		return
	}
	response.OK(c, "dates retrieved successfully", personDates)
}

func (h *PersonDateHandler) UpdatePersonDates(c *gin.Context) {
	var input services.UpdatePersonDatesInput
	userID := c.GetInt("userID")
	if userID == 0 {
		response.Unauthorized(c, "user not authenticated")
		return
	}
	personID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid person ID", nil)
		return
	}
	personDatesID, err := uuid.Parse(c.Param("dateId"))
	if err != nil {
		response.BadRequest(c, "invalid dates ID", nil)
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "validation error", err.Error())
		return
	}

	personDates, err := h.service.UpdatePersonDates(personID, personDatesID, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "dates not found")
			return
		}
		response.InternalError(c, "failed to update dates")
		return
	}
	response.OK(c, "dates updated successfully", personDates)
}

func (h *PersonDateHandler) DeletePersonDates(c *gin.Context) {
	parsedID, err := uuid.Parse(c.Param("dateId"))
	if err != nil {
		response.BadRequest(c, "invalid date ID", nil)
		return
	}
	parsedPersonID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid person ID", nil)
		return
	}

	if err := h.service.DeletePersonDates(parsedPersonID, parsedID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "dates not found")
			return
		}
		response.InternalError(c, "failed to delete dates")
		return
	}
	response.OK(c, "dates deleted successfully", nil)
}