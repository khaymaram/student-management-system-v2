package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"sms-backend/internal/helpers"
	"sms-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

type MajorHandler struct{ repository repositories.MajorRepository }

func NewMajorHandler(repository repositories.MajorRepository) *MajorHandler {
	return &MajorHandler{repository: repository}
}

func (h *MajorHandler) Create(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		helpers.ErrorResponse(c, http.StatusBadRequest, "major name is required")
		return
	}
	major, err := h.repository.Create(req.Name)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	helpers.SuccessResponse(c, http.StatusCreated, major)
}

func (h *MajorHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "invalid major id")
		return
	}
	if err := h.repository.Delete(id); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	helpers.SuccessResponse(c, http.StatusOK, "major deleted")
}

func (h *MajorHandler) GetAll(c *gin.Context) {
	pageParam := c.Query("page")
	limitParam := c.Query("limit")
	if pageParam != "" || limitParam != "" {
		page, limit := 1, 5
		var err error
		if pageParam != "" {
			page, err = strconv.Atoi(pageParam)
			if err != nil || page < 1 {
				helpers.ErrorResponse(c, http.StatusBadRequest, "page must be a positive integer")
				return
			}
		}
		if limitParam != "" {
			limit, err = strconv.Atoi(limitParam)
			if err != nil || limit < 1 || limit > 100 {
				helpers.ErrorResponse(c, http.StatusBadRequest, "limit must be between 1 and 100")
				return
			}
		}
		majors, total, err := h.repository.GetPaginated(page, limit)
		if err != nil {
			helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		totalPages := (total + int64(limit) - 1) / int64(limit)
		helpers.SuccessResponse(c, http.StatusOK, gin.H{
			"data": majors, "page": page, "pageSize": limit, "total": total, "totalPages": totalPages,
		})
		return
	}

	majors, err := h.repository.GetAll()
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.SuccessResponse(c, http.StatusOK, majors)
}
