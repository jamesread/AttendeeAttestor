package main

import (
	"encoding/json"
	"os"
	"sync"
)

type IssuedTicket struct {
	EventID    string                 `json:"eventID"`
	Data       map[string]interface{} `json:"data"`
	Signature  string                 `json:"signature"`
	QRCode     string                 `json:"qrCode"`
}

type TicketStore struct {
	tickets []IssuedTicket
	mutex   sync.RWMutex
}

var globalTicketStore = &TicketStore{
	tickets: make([]IssuedTicket, 0),
}

func getTicketsFilePath() string {
	filePath := os.Getenv("TICKETS_FILE")
	if filePath != "" {
		return filePath
	}
	return "issued-tickets.json"
}

func saveTicketsToFile(filePath string) error {
	globalTicketStore.mutex.RLock()
	tickets := make([]IssuedTicket, len(globalTicketStore.tickets))
	copy(tickets, globalTicketStore.tickets)
	globalTicketStore.mutex.RUnlock()

	fileData, err := json.MarshalIndent(tickets, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, fileData, 0644)
}

func loadTicketsFromFile(filePath string) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var tickets []IssuedTicket
	if err := json.Unmarshal(fileData, &tickets); err != nil {
		return err
	}

	globalTicketStore.mutex.Lock()
	defer globalTicketStore.mutex.Unlock()

	globalTicketStore.tickets = tickets

	return nil
}

func addIssuedTicket(ticket IssuedTicket) error {
	globalTicketStore.mutex.Lock()
	globalTicketStore.tickets = append(globalTicketStore.tickets, ticket)
	tickets := make([]IssuedTicket, len(globalTicketStore.tickets))
	copy(tickets, globalTicketStore.tickets)
	globalTicketStore.mutex.Unlock()

	filePath := getTicketsFilePath()
	fileData, err := json.MarshalIndent(tickets, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, fileData, 0644)
}

func getAllIssuedTickets() []IssuedTicket {
	globalTicketStore.mutex.RLock()
	defer globalTicketStore.mutex.RUnlock()

	tickets := make([]IssuedTicket, len(globalTicketStore.tickets))
	copy(tickets, globalTicketStore.tickets)
	return tickets
}

func getIssuedTicketsForEvent(eventID string) []IssuedTicket {
	globalTicketStore.mutex.RLock()
	defer globalTicketStore.mutex.RUnlock()

	var eventTickets []IssuedTicket
	for _, ticket := range globalTicketStore.tickets {
		if ticket.EventID == eventID {
			eventTickets = append(eventTickets, ticket)
		}
	}
	return eventTickets
}

func getIssuedTicketByID(ticketID string) *IssuedTicket {
	globalTicketStore.mutex.RLock()
	defer globalTicketStore.mutex.RUnlock()

	for _, ticket := range globalTicketStore.tickets {
		ticketIDValue := ticket.Data["ticketId"]
		if ticketIDValue == nil {
			ticketIDValue = ticket.Data["ticketID"]
		}
		if ticketIDValue != nil {
			if ticketIDString, ok := ticketIDValue.(string); ok && ticketIDString == ticketID {
				ticketCopy := ticket
				return &ticketCopy
			}
		}
	}
	return nil
}

