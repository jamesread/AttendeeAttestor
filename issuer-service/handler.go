package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type GenerateResponse struct {
	QRCode    string                 `json:"qrCode"`
	Data      map[string]interface{} `json:"data"`
	Signature string                 `json:"signature"`
}

func createGenerateHandler(defaultPrivateKey ed25519.PrivateKey) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		requestData, err := parseRequestData(request.Body)
		if err != nil {
			http.Error(writer, "Invalid request data", http.StatusBadRequest)
			return
		}

		eventID := ""
		if eventIDValue, exists := requestData["eventID"]; exists {
			if eventIDString, ok := eventIDValue.(string); ok {
				eventID = eventIDString
			}
		}

		var eventConfig *EventConfig
		if eventID != "" {
			config, err := getEventConfig(eventID)
			if err == nil {
				eventConfig = config
				if config.ValidFrom != "" {
					requestData["validFrom"] = convertToISO8601UTC(config.ValidFrom)
				}
				if config.ValidUntil != "" {
					requestData["validUntil"] = convertToISO8601UTC(config.ValidUntil)
				}
			}
		}

		ticketID := ""
		if ticketIDValue, exists := requestData["ticketID"]; exists {
			if ticketIDString, ok := ticketIDValue.(string); ok && ticketIDString != "" {
				ticketID = ticketIDString
			}
		}
		if ticketID == "" {
			ticketID = generateTicketID()
		}
		requestData["ticketId"] = ticketID
		delete(requestData, "ticketID")

		requestData["v"] = 1

		privateKey := getPrivateKeyForRequest(requestData, defaultPrivateKey)

		keyID := ""
		if eventID != "" {
			keyID = eventID
		} else if eventConfig != nil {
			keyID = eventConfig.PublicKey
		} else {
			publicKey := privateKey.Public().(ed25519.PublicKey)
			keyID = base64.StdEncoding.EncodeToString(publicKey)
		}

		coseBytes, err := signCOSE(requestData, privateKey, keyID)
		if err != nil {
			http.Error(writer, "Failed to create COSE signature", http.StatusInternalServerError)
			return
		}

		coseBase64 := base64.RawURLEncoding.EncodeToString(coseBytes)

		qrCode, err := encodeCOSEAsQRCode(coseBase64)
		if err != nil {
			http.Error(writer, "Failed to generate QR code", http.StatusInternalServerError)
			return
		}

		signatureBase64 := base64.StdEncoding.EncodeToString(coseBytes)

		response := GenerateResponse{
			QRCode:    qrCode,
			Data:      requestData,
			Signature: signatureBase64,
		}

		ticket := IssuedTicket{
			EventID:   eventID,
			Data:      requestData,
			Signature: signatureBase64,
			QRCode:    qrCode,
		}

		addIssuedTicket(ticket)

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(response)
	})
}

func getPrivateKeyForRequest(requestData map[string]interface{}, defaultPrivateKey ed25519.PrivateKey) ed25519.PrivateKey {
	eventID, exists := requestData["eventID"]
	if !exists {
		return defaultPrivateKey
	}

	eventIDString, ok := eventID.(string)
	if !ok {
		return defaultPrivateKey
	}

	privateKey, err := getPrivateKeyForEvent(eventIDString)
	if err != nil {
		return defaultPrivateKey
	}

	return privateKey
}

func parseRequestData(body io.Reader) (map[string]interface{}, error) {
	var requestData map[string]interface{}
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&requestData)
	if err != nil {
		return nil, err
	}
	return requestData, nil
}

func generateTicketID() string {
	return uuid.New().String()
}

func convertToISO8601UTC(dateString string) string {
	if dateString == "" {
		return ""
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateString); err == nil {
			return t.UTC().Format(time.RFC3339)
		}
	}

	if t, err := time.Parse("2006-01-02T15:04", dateString); err == nil {
		return t.UTC().Format(time.RFC3339)
	}

	return dateString
}

