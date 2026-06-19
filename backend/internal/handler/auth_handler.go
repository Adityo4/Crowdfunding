package handler

import (
	"net/http"

	"backend/internal/models"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "Input data tidak valid",
			},
		})
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "EMAIL_EXISTS",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "Format email atau password tidak valid",
			},
		})
		return
	}

	res, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
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

	profile, err := h.authService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Profil tidak ditemukan",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": profile,
	})
}

func (h *AuthHandler) GetUsers(c *gin.Context) {
	users, err := h.authService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Gagal memuat pengguna",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
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

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "Format data profil tidak valid: " + err.Error(),
			},
		})
		return
	}

	user, err := h.authService.UpdateProfile(userID, req)
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
		"data": user,
	})
}
