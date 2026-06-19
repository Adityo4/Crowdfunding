package handler

import (
	"net/http"

	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ArticleHandler struct {
	articleService service.ArticleService
}

func NewArticleHandler(articleService service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService}
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
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

	title := c.PostForm("title")
	content := c.PostForm("content")
	isPublished := c.PostForm("isPublished") == "true"

	if title == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "Judul dan Konten Artikel wajib diisi",
			},
		})
		return
	}

	form, err := c.MultipartForm()
	var files []service.UploadedFile
	if err == nil && form != nil {
		formFiles := form.File["images"]
		for _, formFile := range formFiles {
			openedFile, err := formFile.Open()
			if err != nil {
				continue
			}
			defer openedFile.Close()

			files = append(files, service.UploadedFile{
				Reader:      openedFile,
				Size:        formFile.Size,
				ContentType: formFile.Header.Get("Content-Type"),
			})
		}
	}

	article, err := h.articleService.CreateArticle(userID, title, content, isPublished, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Gagal membuat artikel baru: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": article,
	})
}

func (h *ArticleHandler) GetArticles(c *gin.Context) {
	all := c.Query("all") == "true"
	articles, err := h.articleService.GetArticles(all)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Gagal mengambil data artikel",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": articles,
	})
}

func (h *ArticleHandler) GetArticleBySlug(c *gin.Context) {
	slug := c.Param("slug")
	article, err := h.articleService.GetArticleBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Artikel tidak ditemukan",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": article,
	})
}
