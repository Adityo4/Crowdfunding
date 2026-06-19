package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type UploadedFile struct {
	Reader      io.Reader
	Size        int64
	ContentType string
}

type ArticleService interface {
	GetArticles(all bool) ([]models.Article, error)
	GetArticleBySlug(slug string) (*models.Article, error)
	CreateArticle(authorID uuid.UUID, title, content string, isPublished bool, files []UploadedFile) (*models.Article, error)
}

type articleService struct {
	articleRepo repository.ArticleRepository
	minioClient *minio.Client
	bucketName  string
}

func NewArticleService(articleRepo repository.ArticleRepository, minioClient *minio.Client) ArticleService {
	return &articleService{
		articleRepo: articleRepo,
		minioClient: minioClient,
		bucketName:  "crowdfunding-assets",
	}
}

func (s *articleService) GetArticles(all bool) ([]models.Article, error) {
	return s.articleRepo.GetAll(all)
}

func (s *articleService) GetArticleBySlug(slug string) (*models.Article, error) {
	return s.articleRepo.GetBySlug(slug)
}

func (s *articleService) CreateArticle(authorID uuid.UUID, title, content string, isPublished bool, files []UploadedFile) (*models.Article, error) {
	articleID := uuid.New()
	slug := s.slugify(title) + "-" + articleID.String()[:8]

	var imageUrls []string
	if len(files) > 0 {
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

		for idx, file := range files {
			objectName := fmt.Sprintf("articles/%s/img-%d-%d", articleID.String(), idx, time.Now().Unix())
			_, err = s.minioClient.PutObject(ctx, s.bucketName, objectName, file.Reader, file.Size, minio.PutObjectOptions{
				ContentType: file.ContentType,
			})
			if err != nil {
				return nil, fmt.Errorf("gagal mengunggah gambar ke MinIO: %w", err)
			}
			url := fmt.Sprintf("http://localhost:9000/%s/%s", s.bucketName, objectName)
			imageUrls = append(imageUrls, url)
		}
	}

	var coverImageUrl string
	var imagesStr string
	if len(imageUrls) > 0 {
		coverImageUrl = imageUrls[0]
		imagesStr = strings.Join(imageUrls, ",")
	} else {
		coverImageUrl = "https://images.unsplash.com/photo-1559027615-cd4628902d85?auto=format&fit=crop&w=800&q=80"
		imagesStr = coverImageUrl
	}

	now := time.Now()
	var publishedAt *time.Time
	if isPublished {
		publishedAt = &now
	}

	article := &models.Article{
		ID:            articleID,
		AuthorID:      authorID,
		Title:         title,
		Slug:          slug,
		Content:       content,
		CoverImageUrl: coverImageUrl,
		Images:        imagesStr,
		IsPublished:   isPublished,
		PublishedAt:   publishedAt,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.articleRepo.Create(article); err != nil {
		return nil, err
	}

	return article, nil
}

func (s *articleService) slugify(title string) string {
	t := strings.ToLower(title)
	reg, _ := regexp.Compile("[^a-z0-9]+")
	t = reg.ReplaceAllString(t, "-")
	return strings.Trim(t, "-")
}
