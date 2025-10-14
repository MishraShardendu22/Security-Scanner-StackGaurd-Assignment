package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIRequest makes an HTTP request with automatic token rotation on rate limit
func APIRequest(url string, method string, body interface{}, maxRetries int) ([]byte, error) {
	tokenManager := GetTokenManager()

	for attempt := 0; attempt < maxRetries; attempt++ {
		token := tokenManager.GetCurrentToken()

		var bodyReader io.Reader
		if body != nil {
			jsonData, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
			bodyReader = bytes.NewBuffer(jsonData)
		}

		req, err := http.NewRequest(method, url, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		defer resp.Body.Close()
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		// Check for rate limit (HTTP 429) or unauthorized (HTTP 401)
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusUnauthorized {
			// Rotate to next token
			tokenManager.RotateToken()
			fmt.Printf("Token limit reached (status: %d). Rotating to next token (index: %d)\n",
				resp.StatusCode, tokenManager.GetCurrentIndex())

			// If we've tried all tokens, wait a bit before retrying
			if attempt == tokenManager.GetTokenCount()-1 {
				time.Sleep(2 * time.Second)
			}
			continue
		}

		// Check for success
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return responseBody, nil
		}

		// Other error
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return nil, fmt.Errorf("max retries (%d) reached, all tokens exhausted", maxRetries)
}

// GetRequest is a convenience wrapper for GET requests
func GetRequest(url string, maxRetries int) ([]byte, error) {
	return APIRequest(url, http.MethodGet, nil, maxRetries)
}

// PostRequest is a convenience wrapper for POST requests
func PostRequest(url string, body interface{}, maxRetries int) ([]byte, error) {
	return APIRequest(url, http.MethodPost, body, maxRetries)
}
