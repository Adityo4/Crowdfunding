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
	"backend/internal/utils"

	"github.com/google/uuid"
)

type DonationService interface {
	CreateDonation(req models.CreateDonationRequest, userID *uuid.UUID) (*models.Donation, error)
	ProcessCallback(req models.CallbackRequest) (*models.Donation, error)
	GetDonations() ([]models.Donation, error)
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

		// Kirim email konfirmasi pembayaran jika user ada
		if donation.User != nil && donation.User.Email != "" {
			subject := "Konfirmasi Donasi Berhasil - Yayasan Peduli Amal Indonesia"
			htmlBody := fmt.Sprintf(`
				<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e2e8f0; border-radius: 12px;">
					<h2 style="color: #16a34a; text-align: center;">Donasi Anda Berhasil Diterima!</h2>
					<p>Halo <strong>%s</strong>,</p>
					<p>Terima kasih atas kebaikan Anda. Kami ingin mengonfirmasi bahwa pembayaran donasi Anda untuk program penggalangan dana telah berhasil diverifikasi oleh sistem kami.</p>
					
					<div style="background-color: #f8fafc; border: 1px solid #f1f5f9; padding: 15px; border-radius: 8px; margin: 20px 0;">
						<table style="width: 100%%; font-size: 14px; border-collapse: collapse;">
							<tr>
								<td style="padding: 6px 0; color: #64748b;">Program Donasi:</td>
								<td style="padding: 6px 0; font-weight: bold; text-align: right;">%s</td>
							</tr>
							<tr>
								<td style="padding: 6px 0; color: #64748b;">Jumlah Donasi:</td>
								<td style="padding: 6px 0; font-weight: bold; text-align: right; color: #16a34a; font-size: 16px;">Rp %s</td>
							</tr>
							<tr>
								<td style="padding: 6px 0; color: #64748b;">Ref Transaksi:</td>
								<td style="padding: 6px 0; font-family: monospace; text-align: right;">%s</td>
							</tr>
							<tr>
								<td style="padding: 6px 0; color: #64748b;">Status:</td>
								<td style="padding: 6px 0; text-align: right;"><span style="background-color: #dcfce7; color: #16a34a; padding: 3px 8px; border-radius: 9999px; font-size: 12px; font-weight: bold;">BERHASIL</span></td>
							</tr>
						</table>
					</div>

					<p style="font-style: italic; color: #64748b; text-align: center; margin: 25px 0;">"Semoga kebaikan yang Anda tanamkan membawa berkah dan kebahagiaan bagi Anda serta penerima manfaat."</p>
					<hr style="border: 0; border-top: 1px solid #e2e8f0; margin: 20px 0;" />
					<p style="font-size: 12px; color: #94a3b8; text-align: center;">Ini adalah email otomatis, mohon tidak membalas email ini.<br>&copy; 2026 Yayasan Peduli Amal Indonesia. All rights reserved.</p>
				</div>
			`, donation.User.FullName, donation.Charity.Title, fmt.Sprintf("%.0f", donation.Amount), donation.PaymentReference)
			
			go func() {
				_ = utils.SendEmail(donation.User.Email, subject, htmlBody)
			}()
		}
	} else if req.Status == "failed" || req.Status == "expire" {
		donation.Status = "failed"
		donation.UpdatedAt = time.Now()
		_ = s.donationRepo.Update(donation)

		// Kirim email notifikasi donasi gagal/expired jika user ada
		if donation.User != nil && donation.User.Email != "" {
			subject := "Transaksi Donasi Gagal - Yayasan Peduli Amal Indonesia"
			htmlBody := fmt.Sprintf(`
				<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e2e8f0; border-radius: 12px;">
					<h2 style="color: #dc2626; text-align: center;">Transaksi Donasi Gagal</h2>
					<p>Halo <strong>%s</strong>,</p>
					<p>Kami menginformasikan bahwa pembayaran donasi Anda untuk program penggalangan dana telah gagal atau kedaluwarsa.</p>
					
					<div style="background-color: #f8fafc; border: 1px solid #f1f5f9; padding: 15px; border-radius: 8px; margin: 20px 0;">
						<table style="width: 100%%; font-size: 14px; border-collapse: collapse;">
							<tr>
								<td style="padding: 6px 0; color: #64748b;">Program Donasi:</td>
								<td style="padding: 6px 0; font-weight: bold; text-align: right;">%s</td>
							</tr>
							<tr>
								<td style="padding: 6px 0; color: #64748b;">Jumlah Donasi:</td>
								<td style="padding: 6px 0; font-weight: bold; text-align: right;">Rp %s</td>
							</tr>
							<tr>
								<td style="padding: 6px 0; color: #64748b;">Status:</td>
								<td style="padding: 6px 0; text-align: right;"><span style="background-color: #fee2e2; color: #dc2626; padding: 3px 8px; border-radius: 9999px; font-size: 12px; font-weight: bold;">GAGAL/EXPIRED</span></td>
							</tr>
						</table>
					</div>

					<p>Anda dapat mencoba melakukan donasi kembali melalui situs kami jika Anda masih ingin berpartisipasi.</p>
					<hr style="border: 0; border-top: 1px solid #e2e8f0; margin: 20px 0;" />
					<p style="font-size: 12px; color: #94a3b8; text-align: center;">Ini adalah email otomatis, mohon tidak membalas email ini.<br>&copy; 2026 Yayasan Peduli Amal Indonesia. All rights reserved.</p>
				</div>
			`, donation.User.FullName, donation.Charity.Title, fmt.Sprintf("%.0f", donation.Amount))
			
			go func() {
				_ = utils.SendEmail(donation.User.Email, subject, htmlBody)
			}()
		}
	}

	return donation, nil
}

func (s *donationService) GetDonations() ([]models.Donation, error) {
	return s.donationRepo.GetAll()
}
