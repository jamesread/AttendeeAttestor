package main

import (
	"encoding/json"
	"net/http"
)

type ListEventsResponse struct {
	Events []EventInfo `json:"events"`
}

type EventInfo struct {
	EventID   string `json:"eventID"`
	PublicKey string `json:"publicKey"`
}

func createListEventsHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		events := getAllEvents()
		response := ListEventsResponse{
			Events: events,
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(response)
	})
}

func getAllEvents() []EventInfo {
	globalEventStore.mutex.RLock()
	defer globalEventStore.mutex.RUnlock()

	events := make([]EventInfo, 0, len(globalEventStore.events))
	for _, config := range globalEventStore.events {
		events = append(events, EventInfo{
			EventID:   config.EventID,
			PublicKey: config.PublicKey,
		})
	}

	return events
}

