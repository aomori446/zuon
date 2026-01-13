package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"github.com/aomori446/zuon/internal"
)

// AuthenticatedRequest performs an HTTP request.
// If it receives a 401 Unauthorized, it attempts to refresh the access token and retry.
func AuthenticatedRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// 1. Attach current access token
	token := fyne.CurrentApp().Preferences().String("auth_token")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 2. If 401, try to refresh
	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close() // Close the original 401 response

		newToken, err := refreshAccessToken()
		if err != nil {
			// Refresh failed (expired or invalid), return SessionExpired
			return nil, internal.ErrSessionExpired
		}

		// 3. Retry with new token
		req2, err := http.NewRequest(method, url, body) // Need new request object because body might be consumed (if not rewindable)
		// Note: if body is io.Reader and consumed, we might have issues.
		// For this app, Unsplash search uses nil body (GET), so it's fine.
		// If we had POSTs, we'd need to buffer the body.
		if err != nil {
			return nil, err
		}
		
		req2.Header.Set("Authorization", "Bearer "+newToken)
		return client.Do(req2)
	}

	return resp, nil
}

func refreshAccessToken() (string, error) {
	refreshToken := fyne.CurrentApp().Preferences().String("refresh_token")
	if refreshToken == "" {
		return "", errors.New("no refresh token")
	}

	reqBody, _ := json.Marshal(map[string]string{
		"refresh_token": refreshToken,
	})

	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/auth/github/refresh", APIBaseURL),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("refresh failed")
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Save new token
	fyne.CurrentApp().Preferences().SetString("auth_token", result.AccessToken)
	return result.AccessToken, nil
}
