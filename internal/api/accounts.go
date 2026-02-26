// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package api

import (
	"encoding/json"
	"fmt"
)

// GetAccount retrieves the current account details.
func (c *Client) GetAccount() (*Account, error) {
	data, err := c.get(c.accountPath("/account"), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get account: %w", err)
	}

	var account Account
	if err := json.Unmarshal(data, &account); err != nil {
		return nil, fmt.Errorf("unable to parse account response: %w", err)
	}

	return &account, nil
}

// UpdateAccount updates the current account.
func (c *Client) UpdateAccount(account *Account) (*Account, error) {
	data, err := c.put(c.accountPath("/account"), account)
	if err != nil {
		return nil, fmt.Errorf("unable to update account: %w", err)
	}

	var updated Account
	if err := json.Unmarshal(data, &updated); err != nil {
		return nil, fmt.Errorf("unable to parse account response: %w", err)
	}

	return &updated, nil
}

// GetBilling retrieves the billing information for the current account.
func (c *Client) GetBilling() (*Billing, error) {
	data, err := c.get(c.accountPath("/account/billing"), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get billing: %w", err)
	}

	var billing Billing
	if err := json.Unmarshal(data, &billing); err != nil {
		return nil, fmt.Errorf("unable to parse billing response: %w", err)
	}

	return &billing, nil
}
