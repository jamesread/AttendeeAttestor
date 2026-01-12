package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListEventsEndpointReturnsListOfEvents(t *testing.T) {
	createTestEvent("test-event-1")
	createTestEvent("test-event-2")

	request := httptest.NewRequest(http.MethodGet, "/list-events", nil)
	recorder := httptest.NewRecorder()

	handler := createListEventsHandler()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	events, exists := response["events"]
	if !exists {
		t.Fatal("Response should contain events field")
	}

	eventsArray, ok := events.([]interface{})
	if !ok {
		t.Fatal("events should be an array")
	}

	if len(eventsArray) < 2 {
		t.Errorf("Expected at least 2 events, got %d", len(eventsArray))
	}
}

func TestListEventsEndpointReturnsEmptyListWhenNoEvents(t *testing.T) {
	resetEventStore()

	request := httptest.NewRequest(http.MethodGet, "/list-events", nil)
	recorder := httptest.NewRecorder()

	handler := createListEventsHandler()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	events, exists := response["events"]
	if !exists {
		t.Fatal("Response should contain events field")
	}

	eventsArray, ok := events.([]interface{})
	if !ok {
		t.Fatal("events should be an array")
	}

	if len(eventsArray) != 0 {
		t.Errorf("Expected empty events array, got %d events", len(eventsArray))
	}
}

func TestListEventsEndpointReturnsEventDetails(t *testing.T) {
	resetEventStore()
	createTestEvent("test-event-details")

	request := httptest.NewRequest(http.MethodGet, "/list-events", nil)
	recorder := httptest.NewRecorder()

	handler := createListEventsHandler()
	handler.ServeHTTP(recorder, request)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)

	events := response["events"].([]interface{})
	if len(events) == 0 {
		t.Fatal("Expected at least one event")
	}

	event := events[0].(map[string]interface{})

	if _, exists := event["eventID"]; !exists {
		t.Error("Event should contain eventID field")
	}

	if _, exists := event["publicKey"]; !exists {
		t.Error("Event should contain publicKey field")
	}

	eventID := event["eventID"].(string)
	if eventID != "test-event-details" {
		t.Errorf("Expected eventID to be test-event-details, got %s", eventID)
	}
}

func createTestEvent(eventID string) {
	getOrCreateEventConfig(eventID)
}

func resetEventStore() {
	globalEventStore.mutex.Lock()
	defer globalEventStore.mutex.Unlock()
	globalEventStore.events = make(map[string]*EventConfig)
}

