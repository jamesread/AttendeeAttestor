package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveEventsToFileCreatesFile(t *testing.T) {
	testFile := createTestFilePath("test-save.json")
	defer os.Remove(testFile)

	resetEventStore()
	createTestEventConfig("test-event-1", testFile)
	createTestEventConfig("test-event-2", testFile)

	err := saveEventsToFile(testFile)
	if err != nil {
		t.Fatalf("Failed to save events: %v", err)
	}

	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("Event file should be created")
	}
}

func TestSaveEventsToFileContainsEventData(t *testing.T) {
	testFile := createTestFilePath("test-save-data.json")
	defer os.Remove(testFile)

	resetEventStore()
	createTestEventConfig("test-event-save", testFile)

	err := saveEventsToFile(testFile)
	if err != nil {
		t.Fatalf("Failed to save events: %v", err)
	}

	fileData, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read event file: %v", err)
	}

	var savedEvents map[string]EventFileData
	if err := json.Unmarshal(fileData, &savedEvents); err != nil {
		t.Fatalf("Failed to parse event file: %v", err)
	}

	if _, exists := savedEvents["test-event-save"]; !exists {
		t.Error("Saved file should contain test-event-save")
	}
}

func TestLoadEventsFromFileLoadsEvents(t *testing.T) {
	testFile := createTestFilePath("test-load.json")
	defer os.Remove(testFile)

	eventData := map[string]EventFileData{
		"test-event-load": {
			EventID:   "test-event-load",
			PublicKey: "test-public-key",
			PrivateKey: "test-private-key",
		},
	}

	saveTestEventFile(testFile, eventData)

	resetEventStore()
	err := loadEventsFromFile(testFile)
	if err != nil {
		t.Fatalf("Failed to load events: %v", err)
	}

	config, exists := globalEventStore.events["test-event-load"]
	if !exists {
		t.Error("Event should be loaded from file")
	}

	if config.EventID != "test-event-load" {
		t.Errorf("Expected eventID test-event-load, got %s", config.EventID)
	}
}

func TestLoadEventsFromFileHandlesMissingFile(t *testing.T) {
	testFile := createTestFilePath("non-existent.json")

	resetEventStore()
	err := loadEventsFromFile(testFile)

	if err != nil {
		t.Errorf("Loading non-existent file should not error, got %v", err)
	}
}

func TestLoadEventsFromFileHandlesEmptyFile(t *testing.T) {
	testFile := createTestFilePath("empty.json")
	defer os.Remove(testFile)

	os.WriteFile(testFile, []byte("{}"), 0644)

	resetEventStore()
	err := loadEventsFromFile(testFile)
	if err != nil {
		t.Fatalf("Loading empty file should not error, got %v", err)
	}

	if len(globalEventStore.events) != 0 {
		t.Error("Empty file should result in empty event store")
	}
}

func createTestFilePath(filename string) string {
	return filepath.Join(os.TempDir(), filename)
}

func createTestEventConfig(eventID string, filePath string) {
	config, _ := getOrCreateEventConfig(eventID)
	if config != nil {
		saveEventsToFile(filePath)
	}
}

func saveTestEventFile(filePath string, events map[string]EventFileData) {
	fileData, _ := json.Marshal(events)
	os.WriteFile(filePath, fileData, 0644)
}

