package unsplash

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
	
	"github.com/aomori446/zuon/internal"
)

const (
	baseURL   = "https://api.unsplash.com"
	searchURL = "/search/photos"
)

type Client struct {
	accessKey  string
	httpClient *http.Client
	lastReq    time.Time
	mu         sync.Mutex
}

type SearchResult struct {
	Total      int     `json:"total"`
	TotalPages int     `json:"total_pages"`
	Results    []Photo `json:"results"`
}

type Photo struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	URLs        URLs   `json:"urls"`
	User        User   `json:"user"`
}

type URLs struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
	Thumb   string `json:"thumb"`
}

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("API key is empty")
	}
	
	return &Client{
		accessKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (c *Client) SearchPhotos(query string, page int, perPage int) (*SearchResult, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if time.Since(c.lastReq) < 1*time.Second {
		time.Sleep(1 * time.Second)
	}
	c.lastReq = time.Now()
	
	params := url.Values{}
	params.Add("query", query)
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("per_page", fmt.Sprintf("%d", perPage))
	params.Add("orientation", "landscape")
	
	reqURL := fmt.Sprintf("%s%s?%s", baseURL, searchURL, params.Encode())
	
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Client-ID "+c.accessKey)
	req.Header.Set("Accept-Version", "v1")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, internal.ErrNetworkIssue
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, internal.ErrInvalidAPIKey
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, internal.ErrRateLimited
	}
	if resp.StatusCode != http.StatusOK {
		return nil, internal.ErrNetworkIssue
	}
	
	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	return &result, nil
}
