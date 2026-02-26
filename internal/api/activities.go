// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package api

import (
	"encoding/json"
	"fmt"
)

// GetActivities retrieves a list of account activities.
func (c *Client) GetActivities(params map[string]string) ([]Activity, error) {
	data, err := c.get(c.accountPath("/activities"), params)
	if err != nil {
		return nil, fmt.Errorf("unable to get activities: %w", err)
	}

	var activities []Activity
	if err := json.Unmarshal(data, &activities); err != nil {
		return nil, fmt.Errorf("unable to parse activities response: %w", err)
	}

	return activities, nil
}
