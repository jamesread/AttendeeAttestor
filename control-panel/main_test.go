package main

import (
	"os"
	"testing"
)

func TestGetPortReturnsDefaultWhenEnvNotSet(t *testing.T) {
	originalPort := os.Getenv("PORT")
	os.Unsetenv("PORT")
	defer os.Setenv("PORT", originalPort)

	port := getPort()
	if port != "3001" {
		t.Errorf("Expected default port 3001, got %s", port)
	}
}

func TestGetPortReturnsEnvValueWhenSet(t *testing.T) {
	originalPort := os.Getenv("PORT")
	os.Setenv("PORT", "8080")
	defer os.Setenv("PORT", originalPort)

	port := getPort()
	if port != "8080" {
		t.Errorf("Expected port 8080, got %s", port)
	}
}

func TestGetFrontendDirReturnsEnvValueWhenSet(t *testing.T) {
	originalDir := os.Getenv("FRONTEND_DIR")
	os.Setenv("FRONTEND_DIR", "/custom/path")
	defer os.Setenv("FRONTEND_DIR", originalDir)

	dir := getFrontendDir()
	if dir != "/custom/path" {
		t.Errorf("Expected /custom/path, got %s", dir)
	}
}

