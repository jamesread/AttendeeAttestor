package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateEventEndpointAcceptsEventID(t *testing.T) {
	requestBody := `{"eventID":"test-event-123"}`
	request := httptest.NewRequest(http.MethodPost, "/create-event", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler := createCreateEventHandler()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}
}

func TestCreateEventEndpointReturnsQRCode(t *testing.T) {
	requestBody := `{"eventID":"test-event-456"}`
	request := httptest.NewRequest(http.MethodPost, "/create-event", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler := createCreateEventHandler()
	handler.ServeHTTP(recorder, request)

	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, exists := response["qrCode"]; !exists {
		t.Error("Response should contain qrCode field")
	}

	qrCode, ok := response["qrCode"].(string)
	if !ok {
		t.Error("qrCode should be a string")
	}

	if qrCode == "" {
		t.Error("qrCode should not be empty")
	}
}

func TestCreateEventEndpointReturnsPublicKey(t *testing.T) {
	requestBody := `{"eventID":"test-event-789"}`
	request := httptest.NewRequest(http.MethodPost, "/create-event", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler := createCreateEventHandler()
	handler.ServeHTTP(recorder, request)

	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, exists := response["publicKey"]; !exists {
		t.Error("Response should contain publicKey field")
	}

	publicKey, ok := response["publicKey"].(string)
	if !ok {
		t.Error("publicKey should be a string")
	}

	if publicKey == "" {
		t.Error("publicKey should not be empty")
	}
}

func TestCreateEventEndpointReturnsEventID(t *testing.T) {
	requestBody := `{"eventID":"test-event-abc"}`
	request := httptest.NewRequest(http.MethodPost, "/create-event", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler := createCreateEventHandler()
	handler.ServeHTTP(recorder, request)

	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, exists := response["eventID"]; !exists {
		t.Error("Response should contain eventID field")
	}

	eventID, ok := response["eventID"].(string)
	if !ok {
		t.Error("eventID should be a string")
	}

	if eventID != "test-event-abc" {
		t.Errorf("Expected eventID to be test-event-abc, got %s", eventID)
	}
}

func TestCreateEventEndpointRequiresEventID(t *testing.T) {
	requestBody := `{}`
	request := httptest.NewRequest(http.MethodPost, "/create-event", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler := createCreateEventHandler()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", recorder.Code)
	}
}

func TestCreateEventEndpointInitializesSigningKey(t *testing.T) {
	requestBody := `{"eventID":"test-event-key"}`
	request := httptest.NewRequest(http.MethodPost, "/create-event", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler := createCreateEventHandler()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", recorder.Code)
	}

	requestBody2 := `{"eventID":"test-event-key"}`
	request2 := httptest.NewRequest(http.MethodPost, "/create-event", strings.NewReader(requestBody2))
	request2.Header.Set("Content-Type", "application/json")
	recorder2 := httptest.NewRecorder()

	handler.ServeHTTP(recorder2, request2)

	var response1 map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response1)

	var response2 map[string]interface{}
	json.Unmarshal(recorder2.Body.Bytes(), &response2)

	publicKey1 := response1["publicKey"].(string)
	publicKey2 := response2["publicKey"].(string)

	if publicKey1 != publicKey2 {
		t.Error("Same event should return same public key")
	}
}

func TestCreateEventQRCodeIsBase64PNG(t *testing.T) {
	requestBody := `{"eventID":"test-event-qr"}`
	request := httptest.NewRequest(http.MethodPost, "/create-event", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler := createCreateEventHandler()
	handler.ServeHTTP(recorder, request)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)

	qrCode := response["qrCode"].(string)

	if qrCode == "" {
		t.Error("QR code should not be empty")
	}

	decoded, err := base64.StdEncoding.DecodeString(qrCode)
	if err != nil {
		t.Errorf("QR code should be valid base64: %v", err)
	}

	if len(decoded) < 4 {
		t.Error("QR code should contain PNG data")
	}

	pngSignature := []byte{0x89, 0x50, 0x4E, 0x47}
	if len(decoded) < 4 {
		t.Fatal("Decoded data too short")
	}
	if decoded[0] != pngSignature[0] || decoded[1] != pngSignature[1] || decoded[2] != pngSignature[2] || decoded[3] != pngSignature[3] {
		t.Error("QR code should be a PNG image")
	}
}

