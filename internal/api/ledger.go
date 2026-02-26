// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package api

import (
	"encoding/json"
	"fmt"
)

// GetLedgers retrieves a paginated list of ledger entries.
func (c *Client) GetLedgers(params map[string]string) ([]Ledger, error) {
	data, err := c.get(c.accountPath("/ledger"), params)
	if err != nil {
		return nil, fmt.Errorf("unable to get ledgers: %w", err)
	}

	var ledgers []Ledger
	if err := json.Unmarshal(data, &ledgers); err != nil {
		return nil, fmt.Errorf("unable to parse ledgers response: %w", err)
	}

	return ledgers, nil
}

// GetLedger retrieves a single ledger entry by ID.
func (c *Client) GetLedger(id uint) (*Ledger, error) {
	path := fmt.Sprintf("/ledger/%d", id)
	data, err := c.get(c.accountPath(path), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get ledger: %w", err)
	}

	var ledger Ledger
	if err := json.Unmarshal(data, &ledger); err != nil {
		return nil, fmt.Errorf("unable to parse ledger response: %w", err)
	}

	return &ledger, nil
}

// CreateLedger creates a new ledger entry.
func (c *Client) CreateLedger(req *LedgerCreateRequest) (*Ledger, error) {
	data, err := c.post(c.accountPath("/ledger"), req)
	if err != nil {
		return nil, fmt.Errorf("unable to create ledger: %w", err)
	}

	var ledger Ledger
	if err := json.Unmarshal(data, &ledger); err != nil {
		return nil, fmt.Errorf("unable to parse ledger response: %w", err)
	}

	return &ledger, nil
}

// UpdateLedger updates an existing ledger entry.
func (c *Client) UpdateLedger(id uint, req *LedgerUpdateRequest) (*Ledger, error) {
	path := fmt.Sprintf("/ledger/%d", id)
	data, err := c.put(c.accountPath(path), req)
	if err != nil {
		return nil, fmt.Errorf("unable to update ledger: %w", err)
	}

	var ledger Ledger
	if err := json.Unmarshal(data, &ledger); err != nil {
		return nil, fmt.Errorf("unable to parse ledger response: %w", err)
	}

	return &ledger, nil
}

// DeleteLedger deletes a ledger entry by ID.
func (c *Client) DeleteLedger(id uint) error {
	path := fmt.Sprintf("/ledger/%d", id)
	_, err := c.delete(c.accountPath(path))
	if err != nil {
		return fmt.Errorf("unable to delete ledger: %w", err)
	}

	return nil
}

// GetLedgerSummary retrieves a summary of ledger data.
func (c *Client) GetLedgerSummary(params map[string]string) (*LedgerSummary, error) {
	data, err := c.get(c.accountPath("/ledger-summary"), params)
	if err != nil {
		return nil, fmt.Errorf("unable to get ledger summary: %w", err)
	}

	var summary LedgerSummary
	if err := json.Unmarshal(data, &summary); err != nil {
		return nil, fmt.Errorf("unable to parse summary response: %w", err)
	}

	return &summary, nil
}

// GetLedgerPL retrieves the profit and loss summary from ledger data.
func (c *Client) GetLedgerPL(params map[string]string) (*PnlReport, error) {
	data, err := c.get(c.accountPath("/ledger-pl-summary"), params)
	if err != nil {
		return nil, fmt.Errorf("unable to get ledger P&L: %w", err)
	}

	var report PnlReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("unable to parse P&L response: %w", err)
	}

	return &report, nil
}
