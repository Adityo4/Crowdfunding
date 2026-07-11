package handler

import (
	"net/http"

	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ArticleCommentHandler struct {
	commentService service.ArticleCommentService
}

func NewArticleCommentHandler(commentService service.ArticleCommentService) *ArticleCommentHandler {
	return &ArticleCommentHandler{commentService}
}

type CreateCommentRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email"`
	Content string `json:"content" binding:"required"`
}

func (h *ArticleCommentHandler) GetComments(c *gin.Context) {
	slug := c.Param("slug")
	comments, err := h.commentService.GetCommentsByArticleSlug(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Gagal mengambil data komentar: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": comments,
	})
}

func (h *ArticleCommentHandler) CreateComment(c *gin.Context) {
	slug := c.Param("slug")
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "Nama dan Konten komentar wajib diisi",
			},
		})
		return
	}

	comment, err := h.commentService.CreateComment(slug, req.Name, req.Email, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Gagal menyimpan komentar: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": comment,
	})
}
