package handler

import (
	"net/http"

	"backend/internal/models"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DonationHandler struct {
	donationService service.DonationService
}

func NewDonationHandler(donationService service.DonationService) *DonationHandler {
	return &DonationHandler{donationService}
}

func (h *DonationHandler) CreateDonation(c *gin.Context) {
	var req models.CreateDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "Input donasi tidak valid",
			},
		})
		return
	}

	var userUUID *uuid.UUID
	userIDStr, exists := c.Get("userID")
	if exists {
		parsed, err := uuid.Parse(userIDStr.(string))
		if err == nil {
			userUUID = &parsed
		}
	}

	donation, err := h.donationService.CreateDonation(req, userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	// Buat response simulasi invoice
	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"donationId":    donation.ID,
			"amount":        donation.Amount,
			"status":        donation.Status,
			"paymentMethod": donation.PaymentMethod,
			"paymentAction": gin.H{
				"qrUrl":    "https://payment-gateway.com/qr/" + donation.PaymentReference,
				"deeplink": donation.PaymentMethod + "://pay?code=" + donation.PaymentReference,
			},
		},
	})
}

func (h *DonationHandler) ProcessCallback(c *gin.Context) {
	var req models.CallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "Input callback tidak valid",
			},
		})
		return
	}

	donation, err := h.donationService.ProcessCallback(req)
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
		"data": gin.H{
			"donationId": donation.ID,
			"status":     donation.Status,
		},
	})
}
