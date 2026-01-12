<template>
  <div class="tickets-list-container">
    <div class="tickets-list-header">
      <h2>Issued Tickets{{ eventID ? ` - ${eventID}` : '' }}</h2>
      <div class="header-buttons">
        <router-link v-if="eventID" :to="`/issue-ticket?eventID=${eventID}`" class="issue-ticket-button">
          Issue Ticket
        </router-link>
        <router-link v-else to="/issue-ticket" class="issue-ticket-button">
          Issue Ticket
        </router-link>
        <button @click="loadTickets" :disabled="isLoading" class="refresh-button">
          {{ isLoading ? 'Loading...' : 'Refresh' }}
        </button>
      </div>
    </div>

    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>

    <div v-if="isLoading && tickets.length === 0" class="loading-message">
      Loading tickets...
    </div>

    <div v-if="!isLoading && tickets.length === 0 && !errorMessage" class="empty-state">
      <p>No tickets have been issued yet.</p>
    </div>

    <div v-if="tickets.length > 0" class="tickets-table">
      <table>
        <thead>
          <tr>
            <th>Event ID</th>
            <th>Ticket ID</th>
            <th>Ticket Type</th>
            <th>Valid From</th>
            <th>Valid Until</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(ticket, index) in tickets" :key="index">
            <td class="event-id">{{ ticket.eventID || 'N/A' }}</td>
            <td class="ticket-id">
              <router-link 
                v-if="getTicketID(ticket) !== 'N/A'"
                :to="`/ticket-details?ticketID=${encodeURIComponent(getTicketID(ticket))}`"
                class="ticket-id-link"
              >
                {{ getTicketID(ticket) }}
              </router-link>
              <span v-else>{{ getTicketID(ticket) }}</span>
            </td>
            <td class="ticket-type">{{ getTicketType(ticket) }}</td>
            <td class="valid-from">{{ formatDateTime(getValidFrom(ticket)) }}</td>
            <td class="valid-until">{{ formatDateTime(getValidUntil(ticket)) }}</td>
            <td>
              <button @click="showTicketQR(ticket)" class="view-qr-button">
                View QR
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="selectedTicket" class="qr-modal-overlay" @click="closeTicketQR">
      <div class="qr-modal" @click.stop>
        <div class="qr-modal-header">
          <h3>Ticket QR Code</h3>
          <button @click="closeTicketQR" class="close-button">×</button>
        </div>
        <div class="qr-modal-content">
          <div class="ticket-info">
            <p><strong>Event ID:</strong> {{ selectedTicket.eventID || 'N/A' }}</p>
            <p v-if="getTicketID(selectedTicket) && getTicketID(selectedTicket) !== 'N/A'"><strong>Ticket ID:</strong> {{ getTicketID(selectedTicket) }}</p>
            <p v-if="getTicketType(selectedTicket)"><strong>Ticket Type:</strong> {{ getTicketType(selectedTicket) }}</p>
          </div>
          <div class="qr-code-display">
            <img v-if="selectedTicket.qrCode" :src="'data:image/png;base64,' + selectedTicket.qrCode" alt="Ticket QR Code" class="qr-code-image" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { listIssuedTickets } from '../services/ticketApi'

export default {
  name: 'IssuedTicketsList',
  data() {
    return {
      tickets: [],
      isLoading: false,
      errorMessage: '',
      selectedTicket: null,
      eventID: ''
    }
  },
  mounted() {
    this.loadEventIDFromRoute()
    this.loadTickets()
  },
  watch: {
    '$route'(to) {
      this.loadEventIDFromRoute()
      this.loadTickets()
    }
  },
  methods: {
    loadEventIDFromRoute() {
      this.eventID = this.$route.query.eventID || ''
    },
    async loadTickets() {
      this.isLoading = true
      this.errorMessage = ''

      try {
        const response = await listIssuedTickets(this.eventID)
        this.tickets = response.tickets || []
      } catch (error) {
        this.errorMessage = `Failed to load tickets: ${error.message}`
        this.tickets = []
      } finally {
        this.isLoading = false
      }
    },
    getTicketID(ticket) {
      return ticket.data?.ticketId || ticket.data?.ticketID || 'N/A'
    },
    getTicketType(ticket) {
      return ticket.data?.ticketType || 'N/A'
    },
    getValidFrom(ticket) {
      return ticket.data?.validFrom || ''
    },
    getValidUntil(ticket) {
      return ticket.data?.validUntil || ''
    },
    formatDateTime(dateTimeString) {
      if (!dateTimeString) return 'N/A'
      try {
        const date = new Date(dateTimeString)
        return date.toLocaleString()
      } catch {
        return dateTimeString
      }
    },
    showTicketQR(ticket) {
      this.selectedTicket = ticket
    },
    closeTicketQR() {
      this.selectedTicket = null
    }
  }
}
</script>

<style scoped>
.tickets-list-container {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.tickets-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.tickets-list-header h2 {
  color: #2c3e50;
  font-size: 1.5rem;
  margin: 0;
}

.header-buttons {
  display: flex;
  gap: 0.75rem;
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

.tickets-table {
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

.ticket-id {
  font-family: 'Courier New', monospace;
  font-size: 0.9rem;
}

.ticket-id-link {
  color: #3498db;
  text-decoration: none;
  transition: color 0.2s;
}

.ticket-id-link:hover {
  color: #2980b9;
  text-decoration: underline;
}

.ticket-type {
  color: #666;
}

.valid-from,
.valid-until {
  font-size: 0.9rem;
  color: #666;
}

.view-qr-button {
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

.view-qr-button:hover {
  background-color: #2980b9;
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

.ticket-info {
  margin-bottom: 1.5rem;
}

.ticket-info p {
  margin: 0.5rem 0;
  color: #333;
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
</style>

