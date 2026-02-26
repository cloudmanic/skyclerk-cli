// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package api

import (
	"encoding/json"
	"fmt"
)

// GetUsers retrieves all users in the current account.
func (c *Client) GetUsers() ([]User, error) {
	data, err := c.get(c.accountPath("/users"), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get users: %w", err)
	}

	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("unable to parse users response: %w", err)
	}

	return users, nil
}

// RemoveUser removes a user from the current account.
func (c *Client) RemoveUser(id uint) error {
	path := fmt.Sprintf("/users/%d", id)
	_, err := c.delete(c.accountPath(path))
	if err != nil {
		return fmt.Errorf("unable to remove user: %w", err)
	}

	return nil
}

// GetInvites retrieves pending invitations for the current account.
func (c *Client) GetInvites() ([]Invite, error) {
	data, err := c.get(c.accountPath("/users/invite"), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get invites: %w", err)
	}

	var invites []Invite
	if err := json.Unmarshal(data, &invites); err != nil {
		return nil, fmt.Errorf("unable to parse invites response: %w", err)
	}

	return invites, nil
}

// CreateInvite sends an invitation to a user to join the current account.
func (c *Client) CreateInvite(req *InviteCreateRequest) (*Invite, error) {
	data, err := c.post(c.accountPath("/users/invite"), req)
	if err != nil {
		return nil, fmt.Errorf("unable to create invite: %w", err)
	}

	var invite Invite
	if err := json.Unmarshal(data, &invite); err != nil {
		return nil, fmt.Errorf("unable to parse invite response: %w", err)
	}

	return &invite, nil
}

// CancelInvite cancels a pending invitation.
func (c *Client) CancelInvite(id uint) error {
	path := fmt.Sprintf("/user-invite/%d", id)
	_, err := c.delete(c.accountPath(path))
	if err != nil {
		return fmt.Errorf("unable to cancel invite: %w", err)
	}

	return nil
}
