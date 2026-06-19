package service

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"backend/internal/models"
	"backend/internal/repository"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type CharityService interface {
	CreateCharity(req models.CreateCharityRequest, userID uuid.UUID, coverReader io.Reader, coverSize int64, coverContentType string) (*models.Charity, error)
	GetCharities(q string, categorySlug string, status string, sort string, page, perPage int) ([]models.Charity, int, int, error)
	GetCharityBySlug(slug string) (*models.Charity, error)
	UpdateCharityStatus(id uuid.UUID, status string) error
	UpdateCharity(id uuid.UUID, req models.CreateCharityRequest) error
}

type charityService struct {
	charityRepo repository.CharityRepository
	minioClient *minio.Client
	bucketName  string
}

func NewCharityService(charityRepo repository.CharityRepository, minioClient *minio.Client) CharityService {
	return &charityService{
		charityRepo: charityRepo,
		minioClient: minioClient,
		bucketName:  "crowdfunding-assets",
	}
}

func (s *charityService) CreateCharity(req models.CreateCharityRequest, userID uuid.UUID, coverReader io.Reader, coverSize int64, coverContentType string) (*models.Charity, error) {
	charityID := uuid.New()
	slug := s.slugify(req.Title) + "-" + charityID.String()[:8]

	categoryUUID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("kategori ID tidak valid: %w", err)
	}

	var coverURL string
	if coverReader != nil {
		ctx := context.Background()
		// Pastikan bucket tersedia
		exists, err := s.minioClient.BucketExists(ctx, s.bucketName)
		if err != nil {
			return nil, fmt.Errorf("gagal memeriksa bucket MinIO: %w", err)
		}
		if !exists {
			err = s.minioClient.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
			if err != nil {
				return nil, fmt.Errorf("gagal membuat bucket MinIO: %w", err)
			}
			// Set bucket policy publik untuk membaca gambar
			policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, s.bucketName)
			err = s.minioClient.SetBucketPolicy(ctx, s.bucketName, policy)
			if err != nil {
				return nil, fmt.Errorf("gagal mengatur policy bucket: %w", err)
			}
		}

		objectName := fmt.Sprintf("charities/%s/cover-%d", charityID.String(), time.Now().Unix())
		_, err = s.minioClient.PutObject(ctx, s.bucketName, objectName, coverReader, coverSize, minio.PutObjectOptions{
			ContentType: coverContentType,
		})
		if err != nil {
			return nil, fmt.Errorf("gagal mengunggah gambar ke MinIO: %w", err)
		}

		// Asumsi port MinIO external adalah 9000, atau localhost:9000 untuk dev
		coverURL = fmt.Sprintf("http://localhost:9000/%s/%s", s.bucketName, objectName)
	}

	charity := &models.Charity{
		ID:               charityID,
		UserID:           userID,
		CategoryID:       categoryUUID,
		Title:            req.Title,
		Slug:             slug,
		Description:      req.Description,
		TargetAmount:     req.TargetAmount,
		CurrentAmount:    0,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
		CoverImageUrl:    coverURL,
		ContactPerson:    req.ContactPerson,
		ContactEmail:     req.ContactEmail,
		ContactPhone:     req.ContactPhone,
		OrganizationName: req.OrganizationName,
		Status:           "pending", // butuh verifikasi admin dahulu
		IsVerified:       false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = s.charityRepo.Create(charity)
	if err != nil {
		return nil, err
	}

	return charity, nil
}

func (s *charityService) GetCharities(q string, categorySlug string, status string, sort string, page, perPage int) ([]models.Charity, int, int, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	charities, total, err := s.charityRepo.GetAll(q, categorySlug, status, sort, offset, perPage)
	if err != nil {
		return nil, 0, 0, err
	}

	return charities, int(total), page, nil
}

func (s *charityService) GetCharityBySlug(slug string) (*models.Charity, error) {
	return s.charityRepo.GetBySlug(slug)
}

func (s *charityService) UpdateCharityStatus(id uuid.UUID, status string) error {
	charity, err := s.charityRepo.GetByID(id)
	if err != nil {
		return err
	}
	charity.Status = status
	charity.UpdatedAt = time.Now()
	return s.charityRepo.Update(charity)
}

func (s *charityService) UpdateCharity(id uuid.UUID, req models.CreateCharityRequest) error {
	charity, err := s.charityRepo.GetByID(id)
	if err != nil {
		return err
	}

	categoryUUID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return fmt.Errorf("kategori ID tidak valid")
	}

	charity.Title = req.Title
	charity.Description = req.Description
	charity.CategoryID = categoryUUID
	charity.TargetAmount = req.TargetAmount
	charity.StartDate = req.StartDate
	charity.EndDate = req.EndDate
	charity.ContactPerson = req.ContactPerson
	charity.ContactEmail = req.ContactEmail
	charity.ContactPhone = req.ContactPhone
	charity.OrganizationName = req.OrganizationName
	charity.UpdatedAt = time.Now()

	return s.charityRepo.Update(charity)
}

func (s *charityService) slugify(title string) string {
	t := strings.ToLower(title)
	reg, _ := regexp.Compile("[^a-z0-9]+")
	t = reg.ReplaceAllString(t, "-")
	return strings.Trim(t, "-")
}
