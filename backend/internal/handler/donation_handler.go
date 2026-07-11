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

	// Buat response dengan Snap Token Midtrans
	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"donationId":    donation.ID,
			"amount":        donation.Amount,
			"status":        donation.Status,
			"paymentMethod": donation.PaymentMethod,
			"snapToken":     donation.SnapToken,
			"snapUrl":       donation.RedirectURL,
			"paymentAction": gin.H{
				"qrUrl":    "https://payment-gateway.com/qr/" + donation.PaymentReference,
				"deeplink": donation.PaymentMethod + "://pay?code=" + donation.PaymentReference,
			},
		},
	})
}

func (h *DonationHandler) ProcessCallback(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "Input callback tidak valid: " + err.Error(),
			},
		})
		return
	}

	var donationID, paymentRef, status string

	// Deteksi format webhook Midtrans atau simulasi internal
	if orderID, exists := body["order_id"].(string); exists {
		donationID = orderID
		paymentRef, _ = body["transaction_id"].(string)
		status, _ = body["transaction_status"].(string)
	} else {
		donationID, _ = body["donationId"].(string)
		paymentRef, _ = body["paymentReference"].(string)
		status, _ = body["status"].(string)
	}

	if donationID == "" || status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "Donation ID dan status pembayaran wajib diisi",
			},
		})
		return
	}

	donation, err := h.donationService.ProcessCallback(models.CallbackRequest{
		DonationID:       donationID,
		PaymentReference: paymentRef,
		Status:           status,
	})
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
