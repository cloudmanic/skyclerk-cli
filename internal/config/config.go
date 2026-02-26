// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ConfigDir is the directory where the config file is stored.
const ConfigDir = ".config/skyclerk"

// ConfigFile is the name of the config file.
const ConfigFile = "config.json"

// Config holds the CLI configuration including auth credentials and defaults.
type Config struct {
	AccessToken      string `json:"access_token"`
	UserID           uint   `json:"user_id"`
	DefaultAccountID uint   `json:"default_account_id"`
	ApiURL           string `json:"api_url"`
	ClientID         string `json:"client_id"`
}

// DefaultApiURL is the default Skyclerk API URL.
const DefaultApiURL = "https://app.skyclerk.com"

// GetConfigDir returns the full path to the config directory.
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to find home directory: %w", err)
	}

	return filepath.Join(home, ConfigDir), nil
}

// GetConfigPath returns the full path to the config file.
func GetConfigPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, ConfigFile), nil
}

// Load reads the config file from disk and returns a Config struct.
func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	return LoadFromPath(path)
}

// LoadFromPath reads the config file from the given path and returns a Config struct.
func LoadFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not logged in. Run 'skyclerk login' first")
		}
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unable to parse config file: %w", err)
	}

	return &cfg, nil
}

// Save writes the Config struct to disk as JSON.
func Save(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	return SaveToPath(cfg, path)
}

// SaveToPath writes the Config struct to the given path as JSON.
func SaveToPath(cfg *Config, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("unable to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("unable to write config file: %w", err)
	}

	return nil
}

// Delete removes the config file from disk.
func Delete() error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	return DeleteAtPath(path)
}

// DeleteAtPath removes the config file at the given path.
func DeleteAtPath(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to delete config file: %w", err)
	}

	return nil
}

// MaskString masks a string showing only the first 4 and last 4 characters.
func MaskString(s string) string {
	if len(s) <= 8 {
		return "****"
	}

	return s[:4] + "****" + s[len(s)-4:]
}
