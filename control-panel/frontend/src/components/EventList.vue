<template>
  <div class="event-list-container">
    <div class="event-list-header">
      <h2>Events</h2>
      <div class="header-buttons">
        <router-link to="/create-event" class="create-event-button">
          Create Event
        </router-link>
        <button @click="loadEvents" :disabled="isLoading" class="refresh-button">
          {{ isLoading ? 'Loading...' : 'Refresh' }}
        </button>
      </div>
    </div>

    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>

    <div v-if="isLoading && events.length === 0" class="loading-message">
      Loading events...
    </div>

    <div v-if="!isLoading && events.length === 0 && !errorMessage" class="empty-state">
      <p>No events found. Create your first event to get started.</p>
    </div>

    <div v-if="events.length > 0" class="events-table">
      <table>
        <thead>
          <tr>
            <th>Event ID</th>
            <th>Public Key</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="event in events" :key="event.eventID">
            <td class="event-id">{{ event.eventID }}</td>
            <td class="public-key">{{ event.publicKey }}</td>
            <td>
              <div class="action-buttons">
                <button @click="showScannerSetup(event)" class="scanner-setup-button">
                  Scanner Setup
                </button>
                <router-link :to="`/issued-tickets?eventID=${event.eventID}`" class="view-tickets-button">
                  View Tickets
                </router-link>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="selectedEvent" class="qr-modal-overlay" @click="closeScannerSetup">
      <div class="qr-modal" @click.stop>
        <div class="qr-modal-header">
          <h3>Scanner Setup - {{ selectedEvent.eventID }}</h3>
          <button @click="closeScannerSetup" class="close-button">×</button>
        </div>
        <div class="qr-modal-content">
          <p class="qr-instructions">Scan this QR code with a scanner app to configure it for this event:</p>
          <div class="qr-code-display">
            <img v-if="qrCodeBase64" :src="'data:image/png;base64,' + qrCodeBase64" alt="Scanner Setup QR Code" class="qr-code-image" />
            <div v-if="qrCodeError" class="qr-code-error">{{ qrCodeError }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { listEvents, getScannerQR } from '../services/eventApi'

export default {
  name: 'EventList',
  data() {
    return {
      events: [],
      isLoading: false,
      errorMessage: '',
      selectedEvent: null,
      qrCodeBase64: null,
      qrCodeError: null
    }
  },
  mounted() {
    this.loadEvents()
  },
  methods: {
    async loadEvents() {
      this.isLoading = true
      this.errorMessage = ''

      try {
        const response = await listEvents()
        this.events = response.events || []
      } catch (error) {
        this.errorMessage = `Failed to load events: ${error.message}`
        this.events = []
      } finally {
        this.isLoading = false
      }
    },
    async showScannerSetup(event) {
      this.selectedEvent = event
      this.qrCodeBase64 = null
      this.qrCodeError = null
      await this.loadScannerQR(event.eventID)
    },
    async loadScannerQR(eventID) {
      try {
        const response = await getScannerQR(eventID)
        this.qrCodeBase64 = response.qrCode
      } catch (error) {
        this.qrCodeError = `Failed to load QR code: ${error.message}`
      }
    },
    closeScannerSetup() {
      this.selectedEvent = null
      this.qrCodeBase64 = null
      this.qrCodeError = null
    }
  }
}
</script>

<style scoped>
.event-list-container {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.event-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.header-buttons {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}

.create-event-button {
  padding: 0.5rem 1rem;
  background-color: #27ae60;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  text-decoration: none;
  display: inline-block;
  transition: background-color 0.2s;
}

.create-event-button:hover {
  background-color: #229954;
}

.issue-ticket-button {
  padding: 0.5rem 1rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  text-decoration: none;
  display: inline-block;
  transition: background-color 0.2s;
}

.issue-ticket-button:hover {
  background-color: #2980b9;
}

.event-list-header h2 {
  color: #2c3e50;
  font-size: 1.5rem;
  margin: 0;
}

.refresh-button {
  padding: 0.5rem 1rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.refresh-button:hover:not(:disabled) {
  background-color: #2980b9;
}

.refresh-button:disabled {
  background-color: #95a5a6;
  cursor: not-allowed;
}

.error-message {
  padding: 1rem;
  background-color: #fee;
  border: 1px solid #fcc;
  border-radius: 4px;
  color: #c33;
  margin-bottom: 1rem;
}

.loading-message {
  padding: 2rem;
  text-align: center;
  color: #666;
}

.empty-state {
  padding: 2rem;
  text-align: center;
  color: #666;
}

.events-table {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead {
  background-color: #f5f5f5;
}

th {
  padding: 0.75rem;
  text-align: left;
  font-weight: 600;
  color: #2c3e50;
  border-bottom: 2px solid #ddd;
}

td {
  padding: 0.75rem;
  border-bottom: 1px solid #eee;
  color: #333;
}

tbody tr:hover {
  background-color: #f9f9f9;
}

.event-id {
  font-weight: 500;
  color: #2c3e50;
}

.public-key {
  font-family: 'Courier New', monospace;
  font-size: 0.9rem;
  color: #666;
  word-break: break-all;
}

.scanner-setup-button {
  padding: 0.5rem 1rem;
  background-color: #27ae60;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.scanner-setup-button:hover {
  background-color: #229954;
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.issue-ticket-button {
  padding: 0.5rem 1rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  text-decoration: none;
  display: inline-block;
  transition: background-color 0.2s;
}

.issue-ticket-button:hover {
  background-color: #2980b9;
}

.view-tickets-button {
  padding: 0.5rem 1rem;
  background-color: #9b59b6;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  text-decoration: none;
  display: inline-block;
  transition: background-color 0.2s;
}

.view-tickets-button:hover {
  background-color: #8e44ad;
}

.qr-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.qr-modal {
  background: white;
  border-radius: 8px;
  padding: 0;
  max-width: 500px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.qr-modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #ddd;
}

.qr-modal-header h3 {
  margin: 0;
  color: #2c3e50;
  font-size: 1.25rem;
}

.close-button {
  background: none;
  border: none;
  font-size: 2rem;
  color: #666;
  cursor: pointer;
  padding: 0;
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
  transition: color 0.2s;
}

.close-button:hover {
  color: #2c3e50;
}

.qr-modal-content {
  padding: 1.5rem;
}

.qr-instructions {
  margin-bottom: 1.5rem;
  color: #666;
  text-align: center;
}

.qr-code-display {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 1rem;
  background-color: #f9f9f9;
  border-radius: 4px;
}

.qr-code-image {
  max-width: 100%;
  height: auto;
  display: block;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 1rem;
  background-color: white;
}

.qr-code-error {
  color: #c33;
  text-align: center;
  padding: 1rem;
}
</style>

