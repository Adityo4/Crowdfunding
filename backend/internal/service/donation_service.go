package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"backend/internal/models"
	"backend/internal/repository"

	"github.com/google/uuid"
)

type DonationService interface {
	CreateDonation(req models.CreateDonationRequest, userID *uuid.UUID) (*models.Donation, error)
	ProcessCallback(req models.CallbackRequest) (*models.Donation, error)
}

type donationService struct {
	donationRepo repository.DonationRepository
	charityRepo  repository.CharityRepository
}

func NewDonationService(donationRepo repository.DonationRepository, charityRepo repository.CharityRepository) DonationService {
	return &donationService{donationRepo, charityRepo}
}

type MidtransTransactionDetails struct {
	OrderID     string  `json:"order_id"`
	GrossAmount float64 `json:"gross_amount"`
}

type MidtransRequest struct {
	TransactionDetails MidtransTransactionDetails `json:"transaction_details"`
}

type MidtransResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

func (s *donationService) getMidtransSnapToken(orderID string, amount float64) (string, string, error) {
	midtransReq := MidtransRequest{
		TransactionDetails: MidtransTransactionDetails{
			OrderID:     orderID,
			GrossAmount: amount,
		},
	}

	payload, err := json.Marshal(midtransReq)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", "https://app.sandbox.midtrans.com/snap/v1/transactions", bytes.NewBuffer(payload))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	
	serverKey := "SB-Mid-server-r_sQ4bvvoYrkeot4PdJNa1XL:"
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(serverKey))
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("midtrans error status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var midtransResp MidtransResponse
	if err := json.NewDecoder(resp.Body).Decode(&midtransResp); err != nil {
		return "", "", err
	}

	return midtransResp.Token, midtransResp.RedirectURL, nil
}

func (s *donationService) CreateDonation(req models.CreateDonationRequest, userID *uuid.UUID) (*models.Donation, error) {
	// Verifikasi apakah campaign (charity) ada dan aktif
	charity, err := s.charityRepo.GetByID(req.CharityID)
	if err != nil {
		return nil, errors.New("campaign penggalangan dana tidak ditemukan")
	}

	if charity.Status != "active" {
		return nil, errors.New("campaign penggalangan dana tidak aktif")
	}

	donation := &models.Donation{
		ID:            uuid.New(),
		CharityID:     req.CharityID,
		UserID:        userID,
		Amount:        req.Amount,
		Status:        "pending",
		PaymentMethod: req.PaymentMethod,
		Anonymous:     req.Anonymous,
		Message:       req.Message,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Buat reference pembayaran
	donation.PaymentReference = "PAY-REF-" + donation.ID.String()[:8]

	// Dapatkan Snap Token dari Midtrans
	token, redirectURL, err := s.getMidtransSnapToken(donation.ID.String(), donation.Amount)
	if err != nil {
		return nil, fmt.Errorf("gagal menghubungkan ke payment gateway Midtrans: %w", err)
	}
	donation.SnapToken = token
	donation.RedirectURL = redirectURL

	err = s.donationRepo.Create(donation)
	if err != nil {
		return nil, err
	}

	return donation, nil
}

func (s *donationService) ProcessCallback(req models.CallbackRequest) (*models.Donation, error) {
	donationUUID, err := uuid.Parse(req.DonationID)
	if err != nil {
		return nil, errors.New("donation ID tidak valid")
	}

	donation, err := s.donationRepo.GetByID(donationUUID)
	if err != nil {
		return nil, errors.New("transaksi donasi tidak ditemukan")
	}

	if donation.Status == "paid" {
		return donation, nil // sudah lunas
	}

	if req.Status == "settlement" || req.Status == "success" {
		donation.Status = "paid"
		donation.PaymentReference = req.PaymentReference
		donation.UpdatedAt = time.Now()

		err = s.donationRepo.Update(donation)
		if err != nil {
			return nil, err
		}

		// Tambahkan total dana terkumpul di tabel charities
		charity, err := s.charityRepo.GetByID(donation.CharityID)
		if err == nil {
			charity.CurrentAmount += donation.Amount
			_ = s.charityRepo.Update(charity)
		}
	} else if req.Status == "failed" || req.Status == "expire" {
		donation.Status = "failed"
		donation.UpdatedAt = time.Now()
		_ = s.donationRepo.Update(donation)
	}

	return donation, nil
}
