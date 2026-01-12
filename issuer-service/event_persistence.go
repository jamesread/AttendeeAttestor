package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"os"
)

type EventFileData struct {
	EventID    string `json:"eventID"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
	ValidFrom  string `json:"validFrom,omitempty"`
	ValidUntil string `json:"validUntil,omitempty"`
}

func getEventsFilePath() string {
	filePath := os.Getenv("EVENTS_FILE")
	if filePath != "" {
		return filePath
	}
	return "events.json"
}

func saveEventsToFile(filePath string) error {
	globalEventStore.mutex.RLock()
	defer globalEventStore.mutex.RUnlock()

	eventsData := make(map[string]EventFileData)
	for _, config := range globalEventStore.events {
		privateKeyBase64 := base64.StdEncoding.EncodeToString(config.PrivateKey)
		eventsData[config.EventID] = EventFileData{
			EventID:    config.EventID,
			PublicKey:  config.PublicKey,
			PrivateKey: privateKeyBase64,
			ValidFrom:  config.ValidFrom,
			ValidUntil: config.ValidUntil,
		}
	}

	fileData, err := json.MarshalIndent(eventsData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, fileData, 0600)
}

func loadEventsFromFile(filePath string) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var eventsData map[string]EventFileData
	if err := json.Unmarshal(fileData, &eventsData); err != nil {
		return err
	}

	globalEventStore.mutex.Lock()
	defer globalEventStore.mutex.Unlock()

	for _, eventData := range eventsData {
		privateKeyBytes, err := base64.StdEncoding.DecodeString(eventData.PrivateKey)
		if err != nil {
			continue
		}

		privateKey := ed25519.PrivateKey(privateKeyBytes)
		publicKey := privateKey.Public().(ed25519.PublicKey)
		publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKey)

		config := &EventConfig{
			EventID:    eventData.EventID,
			PublicKey:  publicKeyBase64,
			PrivateKey: privateKey,
			ValidFrom:  eventData.ValidFrom,
			ValidUntil: eventData.ValidUntil,
		}

		globalEventStore.events[eventData.EventID] = config
	}

	return nil
}

