package api

import (
	"fmt"
	"net/http"
	"sync"
	
	"github.com/aomori446/zuon/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	// pendingLogins maps request_id to JWT token
	pendingLogins sync.Map
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// Login GET /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	requestID := uuid.New().String()
	// Initialize with empty string meaning "not logged in yet"
	h.pendingLogins.Store(requestID, "")
	
	loginURL := fmt.Sprintf("http://localhost:8080/auth/mock_github?req_id=%s", requestID)
	
	c.JSON(http.StatusOK, gin.H{
		"request_id": requestID,
		"login_url":  loginURL,
	})
}

// MockGitHubPage GET /auth/mock_github
func (h *AuthHandler) MockGitHubPage(c *gin.Context) {
	reqID := c.Query("req_id")
	if reqID == "" {
		c.String(http.StatusBadRequest, "Missing req_id")
		return
	}
	
	html := fmt.Sprintf(`
		<html>
			<body>
				<h1>Mock GitHub Login</h1>
				<p>Click below to simulate GitHub authorization</p>
				<form action="/auth/callback" method="get">
					<input type="hidden" name="req_id" value="%s">
					<button type="submit" style="padding: 10px 20px; font-size: 16px;">Authorize Zuon</button>
				</form>
			</body>
		</html>
	`, reqID)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GET /auth/callback
func (h *AuthHandler) Callback(c *gin.Context) {
	reqID := c.Query("req_id")
	if _, ok := h.pendingLogins.Load(reqID); !ok {
		c.String(http.StatusNotFound, "Invalid or expired request ID")
		return
	}
	
	// Mock user data and generate JWT
	token, err := auth.GenerateToken(123, "github_mock_user")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate token")
		return
	}
	
	// Store the token
	h.pendingLogins.Store(reqID, token)
	
	html := `
		<html>
			<body>
				<h1>Success!</h1>
				<p>You have logged in successfully. You can now close this window and return to the application.</p>
			</body>
		</html>
	`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GET /auth/poll
func (h *AuthHandler) Poll(c *gin.Context) {
	reqID := c.Query("req_id")
	val, ok := h.pendingLogins.Load(reqID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "invalid request_id"})
		return
	}
	
	token := val.(string)
	if token == "" {
		c.JSON(http.StatusOK, gin.H{"status": "pending"})
		return
	}
	
	// Token found, return it and clear the map
	h.pendingLogins.Delete(reqID)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
	})
}
