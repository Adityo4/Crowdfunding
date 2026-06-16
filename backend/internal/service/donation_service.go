package service

import (
	"errors"
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

	// Buat reference pembayaran dummy
	donation.PaymentReference = "PAY-REF-" + donation.ID.String()[:8]

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
