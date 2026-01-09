package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/aomori446/zuon/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	// pendingLogins maps state (request_id) to JWT token
	pendingLogins sync.Map
	clientID      string
	clientSecret  string
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		clientID:     os.Getenv("GITHUB_CLIENT_ID"),
		clientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	}
}

// GET /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	if h.clientID == "" || h.clientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub credentials not configured"})
		return
	}

	state := uuid.New().String()
	// Initialize with empty string meaning "not logged in yet"
	h.pendingLogins.Store(state, "")

	baseURL := os.Getenv("APP_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	redirectURL := fmt.Sprintf("%s/api/v1/auth/github/callback", baseURL)
	
	loginURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user&state=%s",
		h.clientID, redirectURL, state)
	
	c.JSON(http.StatusOK, gin.H{
		"request_id": state,
		"login_url":  loginURL,
	})
}

// GET /auth/callback
func (h *AuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.String(http.StatusBadRequest, "Missing code or state")
		return
	}

	if _, ok := h.pendingLogins.Load(state); !ok {
		c.String(http.StatusBadRequest, "Invalid state")
		return
	}

	// Reconstruct redirect_uri for token exchange
	baseURL := os.Getenv("APP_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	redirectURL := fmt.Sprintf("%s/api/v1/auth/github/callback", baseURL)

	// Exchange code for token
	tokenURL := "https://github.com/login/oauth/access_token"
	payload := map[string]string{
		"client_id":     h.clientID,
		"client_secret": h.clientSecret,
		"code":          code,
		"redirect_uri":  redirectURL,
	}
	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", tokenURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to exchange token")
		return
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse token response")
		return
	}

	if tokenResp.Error != "" {
		c.String(http.StatusBadRequest, "GitHub error: "+tokenResp.Error)
		return
	}

	// Get User Info
	userReq, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	userReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	userResp, err := client.Do(userReq)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to get user info")
		return
	}
	defer userResp.Body.Close()

	var userRespData struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&userRespData); err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse user info")
		return
	}

	// Generate JWT
	jwtToken, err := auth.GenerateToken(userRespData.ID, userRespData.Login)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate JWT")
		return
	}

	// Store the token
	h.pendingLogins.Store(state, jwtToken)

	html := fmt.Sprintf(`
		<html>
			<body style="font-family: sans-serif; text-align: center; padding-top: 50px;">
				<h1 style="color: green;">Success!</h1>
				<p>Hello, <strong>%s</strong>!</p>
				<p>You have logged in successfully.</p>
				<p>You can now close this window and return to Zuon.</p>
				<script>window.close();</script>
			</body>
		</html>
	`, userRespData.Login)
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