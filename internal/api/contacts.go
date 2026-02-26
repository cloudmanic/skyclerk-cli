// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package api

import (
	"encoding/json"
	"fmt"
)

// GetContacts retrieves all contacts for the current account.
func (c *Client) GetContacts(params map[string]string) ([]Contact, error) {
	data, err := c.get(c.accountPath("/contacts"), params)
	if err != nil {
		return nil, fmt.Errorf("unable to get contacts: %w", err)
	}

	var contacts []Contact
	if err := json.Unmarshal(data, &contacts); err != nil {
		return nil, fmt.Errorf("unable to parse contacts response: %w", err)
	}

	return contacts, nil
}

// GetContact retrieves a single contact by ID.
func (c *Client) GetContact(id uint) (*Contact, error) {
	path := fmt.Sprintf("/contacts/%d", id)
	data, err := c.get(c.accountPath(path), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get contact: %w", err)
	}

	var contact Contact
	if err := json.Unmarshal(data, &contact); err != nil {
		return nil, fmt.Errorf("unable to parse contact response: %w", err)
	}

	return &contact, nil
}

// CreateContact creates a new contact.
func (c *Client) CreateContact(req *ContactCreateRequest) (*Contact, error) {
	data, err := c.post(c.accountPath("/contacts"), req)
	if err != nil {
		return nil, fmt.Errorf("unable to create contact: %w", err)
	}

	var contact Contact
	if err := json.Unmarshal(data, &contact); err != nil {
		return nil, fmt.Errorf("unable to parse contact response: %w", err)
	}

	return &contact, nil
}

// UpdateContact updates an existing contact.
func (c *Client) UpdateContact(id uint, req *ContactUpdateRequest) (*Contact, error) {
	path := fmt.Sprintf("/contacts/%d", id)
	data, err := c.put(c.accountPath(path), req)
	if err != nil {
		return nil, fmt.Errorf("unable to update contact: %w", err)
	}

	var contact Contact
	if err := json.Unmarshal(data, &contact); err != nil {
		return nil, fmt.Errorf("unable to parse contact response: %w", err)
	}

	return &contact, nil
}

// DeleteContact deletes a contact by ID.
func (c *Client) DeleteContact(id uint) error {
	path := fmt.Sprintf("/contacts/%d", id)
	_, err := c.delete(c.accountPath(path))
	if err != nil {
		return fmt.Errorf("unable to delete contact: %w", err)
	}

	return nil
}
