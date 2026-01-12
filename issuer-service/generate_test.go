package main

import (
	"crypto/ed25519"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenerateEndpointAcceptsArbitraryKeys(t *testing.T) {
	privateKey, err := generatePrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	requestBody := `{"eventID":"test-event","ticketID":"ticket-123","ticketType":"VIP"}`
	request := httptest.NewRequest(http.MethodPost, "/generate", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler := createGenerateHandler(privateKey)
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	responseBody := recorder.Body.String()
	if !strings.Contains(responseBody, "qrCode") {
		t.Errorf("Response should contain qrCode field")
	}
}

func TestGenerateEndpointSignsJSONContent(t *testing.T) {
	privateKey, err := generatePrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	requestData := map[string]string{
		"eventID":   "test-event",
		"ticketID":  "ticket-123",
		"ticketType": "VIP",
	}
	requestBodyBytes, _ := json.Marshal(requestData)
	request := httptest.NewRequest(http.MethodPost, "/generate", strings.NewReader(string(requestBodyBytes)))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler := createGenerateHandler(privateKey)
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", recorder.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, exists := response["signature"]; !exists {
		t.Error("Response should contain signature field")
	}

	if _, exists := response["data"]; !exists {
		t.Error("Response should contain data field")
	}
}

func TestGenerateEndpointReturnsQRCode(t *testing.T) {
	privateKey, err := generatePrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	requestBody := `{"eventID":"test-event","ticketID":"ticket-123"}`
	request := httptest.NewRequest(http.MethodPost, "/generate", strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler := createGenerateHandler(privateKey)
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", recorder.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	qrCode, exists := response["qrCode"]
	if !exists {
		t.Fatal("Response should contain qrCode field")
	}

	qrCodeString, ok := qrCode.(string)
	if !ok {
		t.Fatal("qrCode should be a string")
	}

	if len(qrCodeString) == 0 {
		t.Error("qrCode should not be empty")
	}
}

func TestSignJSONContentCreatesValidSignature(t *testing.T) {
	privateKey, err := generatePrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	testData := map[string]interface{}{
		"eventID":   "test-event",
		"ticketID":  "ticket-123",
		"ticketType": "VIP",
	}

	signedContent, err := signJSONContent(testData, privateKey)
	if err != nil {
		t.Fatalf("Failed to sign content: %v", err)
	}

	if signedContent.Signature == "" {
		t.Error("Signature should not be empty")
	}

	if signedContent.Data == nil {
		t.Error("Data should not be nil")
	}

	publicKey := privateKey.Public().(ed25519.PublicKey)
	isValid := verifySignature(signedContent.Data, signedContent.Signature, publicKey)
	if !isValid {
		t.Error("Signature should be valid")
	}
}

func TestEncodeSignedContentAsQRCode(t *testing.T) {
	privateKey, err := generatePrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	testData := map[string]interface{}{
		"eventID": "test-event",
		"ticketID": "ticket-123",
	}

	signedContent, err := signJSONContent(testData, privateKey)
	if err != nil {
		t.Fatalf("Failed to sign content: %v", err)
	}

	qrCode, err := encodeSignedContentAsQRCode(signedContent)
	if err != nil {
		t.Fatalf("Failed to encode QR code: %v", err)
	}

	if qrCode == "" {
		t.Error("QR code should not be empty")
	}
}

