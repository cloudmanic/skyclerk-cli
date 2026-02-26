// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package api

import (
	"encoding/json"
	"fmt"
)

// UploadFile uploads a file to the current account, optionally associating it with a ledger entry.
func (c *Client) UploadFile(filePath string, ledgerID string) (*File, error) {
	fields := map[string]string{}
	if ledgerID != "" {
		fields["ledger_id"] = ledgerID
	}

	data, err := c.uploadFile(c.accountPath("/files"), filePath, fields)
	if err != nil {
		return nil, fmt.Errorf("unable to upload file: %w", err)
	}

	var file File
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("unable to parse file response: %w", err)
	}

	return &file, nil
}
