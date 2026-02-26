// Date: 2026-02-25
// Copyright (c) 2026. All rights reserved.

package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestSaveAndLoadConfig verifies that saving and loading a config file round-trips correctly.
func TestSaveAndLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")

	cfg := &Config{
		AccessToken:      "test-token-12345",
		UserID:           42,
		DefaultAccountID: 1,
		ApiURL:           "https://app.skyclerk.com",
	}

	err := SaveToPath(cfg, path)
	if err != nil {
		t.Fatalf("SaveToPath() error = %v", err)
	}

	loaded, err := LoadFromPath(path)
	if err != nil {
		t.Fatalf("LoadFromPath() error = %v", err)
	}

	if loaded.AccessToken != cfg.AccessToken {
		t.Errorf("AccessToken = %q, want %q", loaded.AccessToken, cfg.AccessToken)
	}

	if loaded.UserID != cfg.UserID {
		t.Errorf("UserID = %d, want %d", loaded.UserID, cfg.UserID)
	}

	if loaded.DefaultAccountID != cfg.DefaultAccountID {
		t.Errorf("DefaultAccountID = %d, want %d", loaded.DefaultAccountID, cfg.DefaultAccountID)
	}

	if loaded.ApiURL != cfg.ApiURL {
		t.Errorf("ApiURL = %q, want %q", loaded.ApiURL, cfg.ApiURL)
	}
}

// TestLoadFromPathNotExist verifies that loading a non-existent config file returns a helpful error.
func TestLoadFromPathNotExist(t *testing.T) {
	_, err := LoadFromPath("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("LoadFromPath() expected error for non-existent file, got nil")
	}
}

// TestLoadFromPathInvalidJSON verifies that loading invalid JSON returns an error.
func TestLoadFromPathInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")

	err := os.WriteFile(path, []byte("not json"), 0600)
	if err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err = LoadFromPath(path)
	if err == nil {
		t.Fatal("LoadFromPath() expected error for invalid JSON, got nil")
	}
}

// TestDeleteAtPath verifies that deleting a config file works and is idempotent.
func TestDeleteAtPath(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")

	cfg := &Config{AccessToken: "test"}
	err := SaveToPath(cfg, path)
	if err != nil {
		t.Fatalf("SaveToPath() error = %v", err)
	}

	err = DeleteAtPath(path)
	if err != nil {
		t.Fatalf("DeleteAtPath() error = %v", err)
	}

	// Verify file is gone.
	_, err = os.Stat(path)
	if !os.IsNotExist(err) {
		t.Errorf("expected file to be deleted, but it still exists")
	}

	// Deleting again should not error.
	err = DeleteAtPath(path)
	if err != nil {
		t.Errorf("DeleteAtPath() on missing file error = %v", err)
	}
}

// TestSaveToPathCreatesDirectory verifies that SaveToPath creates parent directories.
func TestSaveToPathCreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "nested", "dir", "config.json")

	cfg := &Config{AccessToken: "test"}
	err := SaveToPath(cfg, path)
	if err != nil {
		t.Fatalf("SaveToPath() error = %v", err)
	}

	loaded, err := LoadFromPath(path)
	if err != nil {
		t.Fatalf("LoadFromPath() error = %v", err)
	}

	if loaded.AccessToken != "test" {
		t.Errorf("AccessToken = %q, want %q", loaded.AccessToken, "test")
	}
}

// TestSaveToPathFilePermissions verifies that the config file is created with 0600 permissions.
func TestSaveToPathFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")

	cfg := &Config{AccessToken: "secret-token"}
	err := SaveToPath(cfg, path)
	if err != nil {
		t.Fatalf("SaveToPath() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("file permissions = %o, want 0600", perm)
	}
}

// TestMaskString verifies string masking for display.
func TestMaskString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"abcdefghijklmnop", "abcd****mnop"},
		{"short", "****"},
		{"12345678", "****"},
		{"123456789", "1234****6789"},
		{"", "****"},
	}

	for _, tt := range tests {
		result := MaskString(tt.input)
		if result != tt.expected {
			t.Errorf("MaskString(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
