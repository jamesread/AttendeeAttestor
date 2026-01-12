package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"sync"
)

type EventConfig struct {
	EventID    string
	PublicKey  string
	PrivateKey ed25519.PrivateKey
	ValidFrom  string
	ValidUntil string
}

type EventStore struct {
	events map[string]*EventConfig
	mutex  sync.RWMutex
}

var globalEventStore = &EventStore{
	events: make(map[string]*EventConfig),
}

func getOrCreateEventConfig(eventID string) (*EventConfig, error) {
	globalEventStore.mutex.RLock()
	config, exists := globalEventStore.events[eventID]
	globalEventStore.mutex.RUnlock()

	if exists {
		return config, nil
	}

	globalEventStore.mutex.Lock()

	config, exists = globalEventStore.events[eventID]
	if exists {
		globalEventStore.mutex.Unlock()
		return config, nil
	}

	privateKey, err := generatePrivateKey()
	if err != nil {
		globalEventStore.mutex.Unlock()
		return nil, err
	}

	publicKey := privateKey.Public().(ed25519.PublicKey)
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKey)

	config = &EventConfig{
		EventID:    eventID,
		PublicKey:  publicKeyBase64,
		PrivateKey: privateKey,
		ValidFrom:  "",
		ValidUntil: "",
	}

	globalEventStore.events[eventID] = config
	globalEventStore.mutex.Unlock()

	filePath := getEventsFilePath()
	saveEventsToFile(filePath)

	return config, nil
}

func getPrivateKeyForEvent(eventID string) (ed25519.PrivateKey, error) {
	config, err := getOrCreateEventConfig(eventID)
	if err != nil {
		return nil, err
	}
	return config.PrivateKey, nil
}

func getEventConfig(eventID string) (*EventConfig, error) {
	return getOrCreateEventConfig(eventID)
}

type ScannerConfig struct {
	EventID   string `json:"eventID"`
	PublicKey string `json:"publicKey"`
}

func createScannerConfig(eventID string, publicKey string) ScannerConfig {
	return ScannerConfig{
		EventID:   eventID,
		PublicKey: publicKey,
	}
}

func encodeScannerConfigAsJSON(scannerConfig ScannerConfig) (string, error) {
	jsonBytes, err := json.Marshal(scannerConfig)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

