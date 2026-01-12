package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type CreateEventRequest struct {
	EventID    string `json:"eventID"`
	ValidFrom  string `json:"validFrom,omitempty"`
	ValidUntil string `json:"validUntil,omitempty"`
}

type CreateEventResponse struct {
	EventID   string `json:"eventID"`
	PublicKey string `json:"publicKey"`
	QRCode    string `json:"qrCode"`
}

func createCreateEventHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		createEventRequest, err := parseCreateEventRequest(request.Body)
		if err != nil {
			http.Error(writer, "Invalid request data", http.StatusBadRequest)
			return
		}

		if createEventRequest.EventID == "" {
			http.Error(writer, "EventID is required", http.StatusBadRequest)
			return
		}

		eventConfig, err := getOrCreateEventConfig(createEventRequest.EventID)
		if err != nil {
			http.Error(writer, "Failed to initialize event", http.StatusInternalServerError)
			return
		}

		if createEventRequest.ValidFrom != "" || createEventRequest.ValidUntil != "" {
			globalEventStore.mutex.Lock()
			if createEventRequest.ValidFrom != "" {
				eventConfig.ValidFrom = createEventRequest.ValidFrom
			}
			if createEventRequest.ValidUntil != "" {
				eventConfig.ValidUntil = createEventRequest.ValidUntil
			}
			globalEventStore.mutex.Unlock()

			filePath := getEventsFilePath()
			saveEventsToFile(filePath)
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

		response := CreateEventResponse{
			EventID:   eventConfig.EventID,
			PublicKey: eventConfig.PublicKey,
			QRCode:    qrCode,
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(response)
	})
}

func parseCreateEventRequest(body io.Reader) (CreateEventRequest, error) {
	var request CreateEventRequest
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&request)
	if err != nil {
		return CreateEventRequest{}, err
	}
	return request, nil
}

