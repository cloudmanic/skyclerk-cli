// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package api

import (
	"encoding/json"
	"fmt"
)

// GetMe retrieves the current user's profile.
func (c *Client) GetMe() (*MeResponse, error) {
	data, err := c.get(c.accountPath("/me"), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get profile: %w", err)
	}

	var me MeResponse
	if err := json.Unmarshal(data, &me); err != nil {
		return nil, fmt.Errorf("unable to parse profile response: %w", err)
	}

	return &me, nil
}

// UpdateMe updates the current user's profile.
func (c *Client) UpdateMe(req *MeUpdateRequest) (*MeResponse, error) {
	data, err := c.put(c.accountPath("/me"), req)
	if err != nil {
		return nil, fmt.Errorf("unable to update profile: %w", err)
	}

	var me MeResponse
	if err := json.Unmarshal(data, &me); err != nil {
		return nil, fmt.Errorf("unable to parse profile response: %w", err)
	}

	return &me, nil
}

// ChangePassword changes the current user's password.
func (c *Client) ChangePassword(req *ChangePasswordRequest) error {
	_, err := c.post(c.accountPath("/me/change-password"), req)
	if err != nil {
		return fmt.Errorf("unable to change password: %w", err)
	}

	return nil
}
