package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Article struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Status   string `json:"status"`
}

var articles = []Article{
	{ID: 1, Title: "Golang Basics", Content: "Introduction to Golang", Category: "Programming", Status: "published"},
}

func main() {
	r := gin.Default()

	r.POST("/article/", createArticle)
	r.GET("/article/:limit/:offset", getArticles)
	r.GET("/article/:id", getArticleByID)
	r.POST("/article/:id", updateArticle)
	r.PUT("/article/:id", updateArticle)
	r.PATCH("/article/:id", updateArticle)
	r.DELETE("/article/:id", deleteArticle)

	r.Run(":8080")
}

// Membuat artikel baru
func createArticle(c *gin.Context) {
	var newArticle Article

	if err := c.ShouldBindJSON(&newArticle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi minimal karakter
	if len(strings.TrimSpace(newArticle.Title)) < 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title minimal 20 karakter"})
		return
	}
	if len(strings.TrimSpace(newArticle.Content)) < 200 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content minimal 200 karakter"})
		return
	}

	newArticle.ID = len(articles) + 1
	articles = append(articles, newArticle)

	c.JSON(http.StatusOK, gin.H{"message": "Article created", "article": newArticle})
}

// Mendapatkan daftar artikel dengan paginasi
func getArticles(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Param("limit"))
	offset, _ := strconv.Atoi(c.Param("offset"))

	if offset >= len(articles) {
		c.JSON(http.StatusOK, []Article{})
		return
	}

	end := offset + limit
	if end > len(articles) {
		end = len(articles)
	}

	c.JSON(http.StatusOK, articles[offset:end])
}

// Mendapatkan artikel berdasarkan ID
func getArticleByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for _, article := range articles {
		if article.ID == id {
			c.JSON(http.StatusOK, article)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
}

// Mengupdate artikel berdasarkan ID
func updateArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedArticle Article

	if err := c.ShouldBindJSON(&updatedArticle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi minimal karakter
	if len(strings.TrimSpace(updatedArticle.Title)) < 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title minimal 20 karakter"})
		return
	}
	if len(strings.TrimSpace(updatedArticle.Content)) < 200 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content minimal 200 karakter"})
		return
	}

	for i, article := range articles {
		if article.ID == id {
			articles[i] = updatedArticle
			articles[i].ID = id
			c.JSON(http.StatusOK, gin.H{"message": "Article updated", "article": articles[i]})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
}

// Menghapus artikel berdasarkan ID
func deleteArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for i, article := range articles {
		if article.ID == id {
			articles = append(articles[:i], articles[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Article deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
}
