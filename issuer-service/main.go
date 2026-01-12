package main

import (
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {
	eventsFilePath := getEventsFilePath()
	if err := loadEventsFromFile(eventsFilePath); err != nil {
		log.Fatalf("Failed to load events from file: %v", err)
	}

	ticketsFilePath := getTicketsFilePath()
	if err := loadTicketsFromFile(ticketsFilePath); err != nil {
		log.Fatalf("Failed to load tickets from file: %v", err)
	}

	defaultPrivateKey, err := loadOrGeneratePrivateKey()
	if err != nil {
		log.Fatalf("Failed to load or generate private key: %v", err)
	}

	generateHandler := createGenerateHandler(defaultPrivateKey)
	corsGenerateHandler := createCORSHandler(generateHandler)
	http.Handle("/generate", corsGenerateHandler)

	createEventHandler := createCreateEventHandler()
	corsCreateEventHandler := createCORSHandler(createEventHandler)
	http.Handle("/create-event", corsCreateEventHandler)

	listEventsHandler := createListEventsHandler()
	corsListEventsHandler := createCORSHandler(listEventsHandler)
	http.Handle("/list-events", corsListEventsHandler)

	getScannerQRHandler := createGetScannerQRHandler()
	corsGetScannerQRHandler := createCORSHandler(getScannerQRHandler)
	http.Handle("/get-scanner-qr/", corsGetScannerQRHandler)

	listIssuedTicketsHandler := createListIssuedTicketsHandler()
	corsListIssuedTicketsHandler := createCORSHandler(listIssuedTicketsHandler)
	http.Handle("/issued-tickets", corsListIssuedTicketsHandler)

	decodeCBORHandler := createDecodeCBORHandler()
	corsDecodeCBORHandler := createCORSHandler(decodeCBORHandler)
	http.Handle("/decode-cbor", corsDecodeCBORHandler)

	getTicketByIDHandler := createGetTicketByIDHandler()
	corsGetTicketByIDHandler := createCORSHandler(getTicketByIDHandler)
	http.Handle("/get-ticket-by-id", corsGetTicketByIDHandler)

	logrus.Info("Starting AttendeeAttestor service on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
