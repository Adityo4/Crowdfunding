package handler

import (
	"net/http"
	"strconv"

	"backend/internal/models"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CharityHandler struct {
	charityService service.CharityService
}

func NewCharityHandler(charityService service.CharityService) *CharityHandler {
	return &CharityHandler{charityService}
}

func (h *CharityHandler) CreateCharity(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Token tidak valid",
			},
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "User ID tidak valid",
			},
		})
		return
	}

	var req models.CreateCharityRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "Form data tidak valid: " + err.Error(),
			},
		})
		return
	}

	// Baca file cover image
	file, header, err := c.Request.FormFile("coverImage")
	var fileReader = file
	var fileSize int64
	var fileContentType string

	if err == nil {
		defer file.Close()
		fileSize = header.Size
		fileContentType = header.Header.Get("Content-Type")
	}

	charity, err := h.charityService.CreateCharity(req, userID, fileReader, fileSize, fileContentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": charity,
	})
}

func (h *CharityHandler) GetCharities(c *gin.Context) {
	q := c.Query("q")
	categorySlug := c.Query("category")
	status := c.Query("status")
	sort := c.Query("sort")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "10"))

	charities, total, currentPage, err := h.charityService.GetCharities(q, categorySlug, status, sort, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Gagal mengambil data campaign",
			},
		})
		return
	}

	totalPages := (total + perPage - 1) / perPage
	hasNext := page < totalPages
	hasPrev := page > 1

	c.JSON(http.StatusOK, gin.H{
		"data": charities,
		"pagination": gin.H{
			"page":       currentPage,
			"perPage":    perPage,
			"total":      total,
			"totalPages": totalPages,
			"hasNext":    hasNext,
			"hasPrev":    hasPrev,
		},
	})
}

func (h *CharityHandler) GetCharityBySlug(c *gin.Context) {
	slug := c.Param("slug")
	charity, err := h.charityService.GetCharityBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Campaign penggalangan dana tidak ditemukan",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": charity,
	})
}

func (h *CharityHandler) UpdateStatus(c *gin.Context) {
	charityIDStr := c.Param("id")
	charityID, err := uuid.Parse(charityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "Charity ID tidak valid",
			},
		})
		return
	}

	var body struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "Body request tidak valid",
			},
		})
		return
	}

	err = h.charityService.UpdateCharityStatus(charityID, body.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status penggalangan dana berhasil diperbarui",
	})
}

func (h *CharityHandler) UpdateCharity(c *gin.Context) {
	charityIDStr := c.Param("id")
	charityID, err := uuid.Parse(charityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "Charity ID tidak valid",
			},
		})
		return
	}

	var req models.CreateCharityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "Data body tidak valid: " + err.Error(),
			},
		})
		return
	}

	err = h.charityService.UpdateCharity(charityID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Penggalangan dana berhasil diperbarui",
	})
}
