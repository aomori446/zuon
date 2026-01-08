package api

import (
	"net/http"
	"strconv"

	"github.com/aomori446/zuon/internal/unsplash"
	"github.com/gin-gonic/gin"
)

type UnsplashHandler struct {
	client *unsplash.Client
}

func NewUnsplashHandler(client *unsplash.Client) *UnsplashHandler {
	return &UnsplashHandler{client: client}
}

func (h *UnsplashHandler) Search(c *gin.Context) {
	query := c.Query("query")
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "12")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil {
		perPage = 12
	}

	results, err := h.client.SearchPhotos(query, page, perPage)
	if err != nil {
		// In a real app, distinguish between client errors and server errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
