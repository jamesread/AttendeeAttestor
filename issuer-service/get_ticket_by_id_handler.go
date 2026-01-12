package main

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type GetTicketByIDResponse struct {
	Ticket *IssuedTicket `json:"ticket,omitempty"`
	Error  string        `json:"error,omitempty"`
}

func createGetTicketByIDHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		logrus.WithFields(logrus.Fields{
			"method": request.Method,
			"path":   request.URL.Path,
		}).Info("Get ticket by ID request received")

		if request.Method != http.MethodGet {
			logrus.Warn("Invalid method for get-ticket-by-id endpoint")
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ticketID := request.URL.Query().Get("ticketID")
		if ticketID == "" {
			logrus.Warn("Missing ticketID parameter")
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			response := GetTicketByIDResponse{
				Error: "ticketID parameter is required",
			}
			json.NewEncoder(writer).Encode(response)
			return
		}

		ticket := getIssuedTicketByID(ticketID)
		if ticket == nil {
			logrus.WithField("ticketID", ticketID).Info("Ticket not found")
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusNotFound)
			response := GetTicketByIDResponse{
				Error: "Ticket not found",
			}
			json.NewEncoder(writer).Encode(response)
			return
		}

		logrus.WithField("ticketID", ticketID).Info("Ticket found")
		response := GetTicketByIDResponse{
			Ticket: ticket,
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(writer).Encode(response); err != nil {
			logrus.WithError(err).Error("Failed to encode response")
		}
	})
}

