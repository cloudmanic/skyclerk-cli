// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package api

import (
	"encoding/json"
	"fmt"
)

// GetLabels retrieves all labels for the current account.
func (c *Client) GetLabels(params map[string]string) ([]Label, error) {
	data, err := c.get(c.accountPath("/labels"), params)
	if err != nil {
		return nil, fmt.Errorf("unable to get labels: %w", err)
	}

	var labels []Label
	if err := json.Unmarshal(data, &labels); err != nil {
		return nil, fmt.Errorf("unable to parse labels response: %w", err)
	}

	return labels, nil
}

// GetLabel retrieves a single label by ID.
func (c *Client) GetLabel(id uint) (*Label, error) {
	path := fmt.Sprintf("/labels/%d", id)
	data, err := c.get(c.accountPath(path), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get label: %w", err)
	}

	var label Label
	if err := json.Unmarshal(data, &label); err != nil {
		return nil, fmt.Errorf("unable to parse label response: %w", err)
	}

	return &label, nil
}

// CreateLabel creates a new label.
func (c *Client) CreateLabel(req *LabelCreateRequest) (*Label, error) {
	data, err := c.post(c.accountPath("/labels"), req)
	if err != nil {
		return nil, fmt.Errorf("unable to create label: %w", err)
	}

	var label Label
	if err := json.Unmarshal(data, &label); err != nil {
		return nil, fmt.Errorf("unable to parse label response: %w", err)
	}

	return &label, nil
}

// UpdateLabel updates an existing label.
func (c *Client) UpdateLabel(id uint, req *LabelUpdateRequest) (*Label, error) {
	path := fmt.Sprintf("/labels/%d", id)
	data, err := c.put(c.accountPath(path), req)
	if err != nil {
		return nil, fmt.Errorf("unable to update label: %w", err)
	}

	var label Label
	if err := json.Unmarshal(data, &label); err != nil {
		return nil, fmt.Errorf("unable to parse label response: %w", err)
	}

	return &label, nil
}

// DeleteLabel deletes a label by ID.
func (c *Client) DeleteLabel(id uint) error {
	path := fmt.Sprintf("/labels/%d", id)
	_, err := c.delete(c.accountPath(path))
	if err != nil {
		return fmt.Errorf("unable to delete label: %w", err)
	}

	return nil
}
