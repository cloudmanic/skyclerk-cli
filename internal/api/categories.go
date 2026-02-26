// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package api

import (
	"encoding/json"
	"fmt"
)

// GetCategories retrieves all categories for the current account.
func (c *Client) GetCategories(params map[string]string) ([]Category, error) {
	data, err := c.get(c.accountPath("/categories"), params)
	if err != nil {
		return nil, fmt.Errorf("unable to get categories: %w", err)
	}

	var categories []Category
	if err := json.Unmarshal(data, &categories); err != nil {
		return nil, fmt.Errorf("unable to parse categories response: %w", err)
	}

	return categories, nil
}

// GetCategory retrieves a single category by ID.
func (c *Client) GetCategory(id uint) (*Category, error) {
	path := fmt.Sprintf("/categories/%d", id)
	data, err := c.get(c.accountPath(path), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get category: %w", err)
	}

	var category Category
	if err := json.Unmarshal(data, &category); err != nil {
		return nil, fmt.Errorf("unable to parse category response: %w", err)
	}

	return &category, nil
}

// CreateCategory creates a new category.
func (c *Client) CreateCategory(req *CategoryCreateRequest) (*Category, error) {
	data, err := c.post(c.accountPath("/categories"), req)
	if err != nil {
		return nil, fmt.Errorf("unable to create category: %w", err)
	}

	var category Category
	if err := json.Unmarshal(data, &category); err != nil {
		return nil, fmt.Errorf("unable to parse category response: %w", err)
	}

	return &category, nil
}

// UpdateCategory updates an existing category.
func (c *Client) UpdateCategory(id uint, req *CategoryUpdateRequest) (*Category, error) {
	path := fmt.Sprintf("/categories/%d", id)
	data, err := c.put(c.accountPath(path), req)
	if err != nil {
		return nil, fmt.Errorf("unable to update category: %w", err)
	}

	var category Category
	if err := json.Unmarshal(data, &category); err != nil {
		return nil, fmt.Errorf("unable to parse category response: %w", err)
	}

	return &category, nil
}

// DeleteCategory deletes a category by ID.
func (c *Client) DeleteCategory(id uint) error {
	path := fmt.Sprintf("/categories/%d", id)
	_, err := c.delete(c.accountPath(path))
	if err != nil {
		return fmt.Errorf("unable to delete category: %w", err)
	}

	return nil
}
