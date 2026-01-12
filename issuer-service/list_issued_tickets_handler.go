package main

import (
	"encoding/json"
	"net/http"
)

type ListIssuedTicketsResponse struct {
	Tickets []IssuedTicket `json:"tickets"`
}

func createListIssuedTicketsHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		eventID := request.URL.Query().Get("eventID")
		var tickets []IssuedTicket

		if eventID != "" {
			tickets = getIssuedTicketsForEvent(eventID)
		} else {
			tickets = getAllIssuedTickets()
		}

		response := ListIssuedTicketsResponse{
			Tickets: tickets,
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(response)
	})
}

