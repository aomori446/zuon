package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/aomori446/zuon/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type loginRequest struct {
	accessToken  string
	refreshToken string
	createdAt    time.Time
}

type AuthHandler struct {
	// pendingLogins maps state (request_id) to loginRequest struct
	pendingLogins sync.Map
	clientID      string
	clientSecret  string
}

func NewAuthHandler() *AuthHandler {
	h := &AuthHandler{
		clientID:     os.Getenv("GITHUB_CLIENT_ID"),
		clientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	}

	// Start cleanup routine
	go h.cleanupLoop()

	return h
}

func (h *AuthHandler) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		h.pendingLogins.Range(func(key, value interface{}) bool {
			req, ok := value.(loginRequest)
			if !ok {
				h.pendingLogins.Delete(key)
				return true
			}
			// Expire requests older than 10 minutes
			if now.Sub(req.createdAt) > 10*time.Minute {
				h.pendingLogins.Delete(key)
			}
			return true
		})
	}
}

// GET /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	if h.clientID == "" || h.clientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub credentials not configured"})
		return
	}

	state := uuid.New().String()
	// Initialize with empty token and current timestamp
	h.pendingLogins.Store(state, loginRequest{accessToken: "", createdAt: time.Now()})

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

	val, ok := h.pendingLogins.Load(state)
	if !ok {
		c.String(http.StatusBadRequest, "Invalid or expired state")
		return
	}
	req := val.(loginRequest)

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

	postReq, _ := http.NewRequest("POST", tokenURL, bytes.NewBuffer(jsonData))
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(postReq)
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

	// Generate Tokens
	accessToken, err := auth.GenerateToken(userRespData.ID, userRespData.Login)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(userRespData.ID, userRespData.Login)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	// Store the tokens
	req.accessToken = accessToken
	req.refreshToken = refreshToken
	h.pendingLogins.Store(state, req)

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
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "invalid or expired request_id"})
		return
	}

	req := val.(loginRequest)
	if req.accessToken == "" {
		c.JSON(http.StatusOK, gin.H{"status": "pending"})
		return
	}

	// Token found, return it and clear the map
	h.pendingLogins.Delete(reqID)
	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"access_token":  req.accessToken,
		"refresh_token": req.refreshToken,
	})
}

// POST /auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	claims, err := auth.ValidateRefreshToken(body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate new access token
	newAccessToken, err := auth.GenerateToken(claims.UserID, claims.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}