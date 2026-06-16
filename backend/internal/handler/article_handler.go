package handler

import (
	"net/http"

	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleService service.ArticleService
}

func NewArticleHandler(articleService service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService}
}

func (h *ArticleHandler) GetArticles(c *gin.Context) {
	articles, err := h.articleService.GetArticles()
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
