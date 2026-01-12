package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetScannerQREndpointReturnsQRCode(t *testing.T) {
	createTestEvent("test-event-qr")

	request := httptest.NewRequest(http.MethodGet, "/get-scanner-qr/test-event-qr", nil)
	recorder := httptest.NewRecorder()

	handler := createGetScannerQRHandler()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

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

func TestGetScannerQREndpointReturnsEventID(t *testing.T) {
	createTestEvent("test-event-id")

	request := httptest.NewRequest(http.MethodGet, "/get-scanner-qr/test-event-id", nil)
	recorder := httptest.NewRecorder()

	handler := createGetScannerQRHandler()
	handler.ServeHTTP(recorder, request)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)

	eventID, exists := response["eventID"]
	if !exists {
		t.Error("Response should contain eventID field")
	}

	if eventID != "test-event-id" {
		t.Errorf("Expected eventID to be test-event-id, got %s", eventID)
	}
}

func TestGetScannerQREndpointReturnsPublicKey(t *testing.T) {
	createTestEvent("test-event-key")

	request := httptest.NewRequest(http.MethodGet, "/get-scanner-qr/test-event-key", nil)
	recorder := httptest.NewRecorder()

	handler := createGetScannerQRHandler()
	handler.ServeHTTP(recorder, request)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)

	if _, exists := response["publicKey"]; !exists {
		t.Error("Response should contain publicKey field")
	}
}

func TestGetScannerQREndpointReturns404ForNonExistentEvent(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/get-scanner-qr/non-existent-event", nil)
	recorder := httptest.NewRecorder()

	handler := createGetScannerQRHandler()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", recorder.Code)
	}
}

func TestGetScannerQRQRCodeIsBase64PNG(t *testing.T) {
	createTestEvent("test-event-png")

	request := httptest.NewRequest(http.MethodGet, "/get-scanner-qr/test-event-png", nil)
	recorder := httptest.NewRecorder()

	handler := createGetScannerQRHandler()
	handler.ServeHTTP(recorder, request)

	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)

	qrCode := response["qrCode"].(string)

	if qrCode == "" {
		t.Fatal("QR code should not be empty")
	}

	decoded, err := base64.StdEncoding.DecodeString(qrCode)
	if err != nil {
		t.Errorf("QR code should be valid base64: %v", err)
	}

	if len(decoded) < 4 {
		t.Fatal("Decoded data too short")
	}

	pngSignature := []byte{0x89, 0x50, 0x4E, 0x47}
	if decoded[0] != pngSignature[0] || decoded[1] != pngSignature[1] || decoded[2] != pngSignature[2] || decoded[3] != pngSignature[3] {
		t.Error("QR code should be a PNG image")
	}
}

