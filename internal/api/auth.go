// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package api

import (
	"encoding/json"
	"fmt"
)

// Login authenticates a user with email, password, and client ID, returning an access token.
func (c *Client) Login(email, password, clientID string) (*LoginResponse, error) {
	reqBody := LoginRequest{
		Username:  email,
		Password:  password,
		GrantType: "password",
		ClientID:  clientID,
	}

	data, err := c.postNoAuth("/oauth/token", reqBody)
	if err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}

	var resp LoginResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unable to parse login response: %w", err)
	}

	return &resp, nil
}

// Logout revokes the current access token.
func (c *Client) Logout() error {
	_, err := c.get("/oauth/logout", map[string]string{
		"access_token": c.accessToken,
	})

	if err != nil {
		return fmt.Errorf("logout failed: %w", err)
	}

	return nil
}

// GetAuthUser returns the currently authenticated user profile.
func (c *Client) GetAuthUser() (*User, error) {
	data, err := c.get("/oauth/me", nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get authenticated user: %w", err)
	}

	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("unable to parse user response: %w", err)
	}

	return &user, nil
}
