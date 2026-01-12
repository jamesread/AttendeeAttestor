package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type GetScannerQRResponse struct {
	EventID   string `json:"eventID"`
	PublicKey string `json:"publicKey"`
	QRCode    string `json:"qrCode"`
}

func createGetScannerQRHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		eventID := extractEventIDFromPath(request.URL.Path)
		if eventID == "" {
			http.Error(writer, "EventID is required", http.StatusBadRequest)
			return
		}

		eventConfig, err := getOrCreateEventConfig(eventID)
		if err != nil {
			http.Error(writer, "Failed to get event config", http.StatusInternalServerError)
			return
		}

		globalEventStore.mutex.RLock()
		_, exists := globalEventStore.events[eventID]
		globalEventStore.mutex.RUnlock()

		if !exists {
			http.Error(writer, "Event not found", http.StatusNotFound)
			return
		}

		scannerConfig := createScannerConfig(eventConfig.EventID, eventConfig.PublicKey)
		scannerConfigJSON, err := encodeScannerConfigAsJSON(scannerConfig)
		if err != nil {
			http.Error(writer, "Failed to encode scanner config", http.StatusInternalServerError)
			return
		}

		qrCode, err := encodeStringAsQRCode(scannerConfigJSON)
		if err != nil {
			http.Error(writer, "Failed to generate QR code", http.StatusInternalServerError)
			return
		}

		response := GetScannerQRResponse{
			EventID:   eventConfig.EventID,
			PublicKey: eventConfig.PublicKey,
			QRCode:    qrCode,
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(response)
	})
}

func extractEventIDFromPath(path string) string {
	prefix := "/get-scanner-qr/"
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	return strings.TrimPrefix(path, prefix)
}

