// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// Client manages HTTP interactions with the Skyclerk API.
type Client struct {
	baseURL    string
	accessToken string
	accountID  uint
	httpClient *http.Client
}

// NewClient creates a new API client with the given access token and base URL.
func NewClient(baseURL string, accessToken string, accountID uint) *Client {
	return &Client{
		baseURL:    baseURL,
		accessToken: accessToken,
		accountID:  accountID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetBaseURL overrides the base URL (useful for testing).
func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

// SetAccountID sets the account ID for API calls.
func (c *Client) SetAccountID(id uint) {
	c.accountID = id
}

// accountPath builds a URL path with the account ID prefix.
func (c *Client) accountPath(path string) string {
	return fmt.Sprintf("/api/v3/%d%s", c.accountID, path)
}

// get performs an authenticated GET request to the given path with query parameters.
func (c *Client) get(path string, params map[string]string) ([]byte, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	q := u.Query()
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Accept", "application/json")

	return c.doRequest(req)
}

// post performs an authenticated POST request with a JSON body.
func (c *Client) post(path string, body interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return c.doRequest(req)
}

// postNoAuth performs a POST request without authentication (for login).
func (c *Client) postNoAuth(path string, body interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return c.doRequest(req)
}

// put performs an authenticated PUT request with a JSON body.
func (c *Client) put(path string, body interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal request body: %w", err)
	}

	req, err := http.NewRequest("PUT", c.baseURL+path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return c.doRequest(req)
}

// delete performs an authenticated DELETE request.
func (c *Client) delete(path string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Accept", "application/json")

	return c.doRequest(req)
}

// uploadFile performs an authenticated multipart file upload.
func (c *Client) uploadFile(path string, filePath string, fields map[string]string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add additional form fields.
	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			return nil, fmt.Errorf("unable to write form field: %w", err)
		}
	}

	// Add the file field.
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("unable to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("unable to copy file data: %w", err)
	}

	writer.Close()

	req, err := http.NewRequest("POST", c.baseURL+path, &buf)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	return c.doRequest(req)
}

// doRequest executes an HTTP request and returns the response body.
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}
