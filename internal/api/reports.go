// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package api

import (
	"encoding/json"
	"fmt"
)

// GetPnlReport retrieves a profit and loss report for the given date range.
func (c *Client) GetPnlReport(params map[string]string) (*PnlReport, error) {
	data, err := c.get(c.accountPath("/reports/pnl"), params)
	if err != nil {
		return nil, fmt.Errorf("unable to get P&L report: %w", err)
	}

	var report PnlReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("unable to parse P&L report response: %w", err)
	}

	return &report, nil
}

// GetPnlByLabel retrieves a P&L report broken down by label.
func (c *Client) GetPnlByLabel(params map[string]string) (*PnlReport, error) {
	data, err := c.get(c.accountPath("/reports/pnl/label"), params)
	if err != nil {
		return nil, fmt.Errorf("unable to get P&L by label: %w", err)
	}

	var report PnlReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("unable to parse P&L by label response: %w", err)
	}

	return &report, nil
}

// GetPnlByCategory retrieves a P&L report broken down by category.
func (c *Client) GetPnlByCategory(params map[string]string) (*PnlReport, error) {
	data, err := c.get(c.accountPath("/reports/pnl/category"), params)
	if err != nil {
		return nil, fmt.Errorf("unable to get P&L by category: %w", err)
	}

	var report PnlReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("unable to parse P&L by category response: %w", err)
	}

	return &report, nil
}

// GetPnlCurrent retrieves the current year P&L report.
func (c *Client) GetPnlCurrent() (*PnlReport, error) {
	data, err := c.get(c.accountPath("/reports/pnl/current"), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get current P&L: %w", err)
	}

	var report PnlReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("unable to parse current P&L response: %w", err)
	}

	return &report, nil
}

// GetIncomeByContact retrieves income broken down by contact.
func (c *Client) GetIncomeByContact(params map[string]string) (*PnlReport, error) {
	data, err := c.get(c.accountPath("/reports/income/by-contact"), params)
	if err != nil {
		return nil, fmt.Errorf("unable to get income by contact: %w", err)
	}

	var report PnlReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("unable to parse income by contact response: %w", err)
	}

	return &report, nil
}

// GetExpensesByContact retrieves expenses broken down by contact.
func (c *Client) GetExpensesByContact(params map[string]string) (*PnlReport, error) {
	data, err := c.get(c.accountPath("/reports/expenses/by-contact"), params)
	if err != nil {
		return nil, fmt.Errorf("unable to get expenses by contact: %w", err)
	}

	var report PnlReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("unable to parse expenses by contact response: %w", err)
	}

	return &report, nil
}
