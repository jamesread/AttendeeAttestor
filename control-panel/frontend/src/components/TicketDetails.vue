<template>
  <div class="ticket-details-container">
    <div class="ticket-details-header">
      <h2>Ticket Details</h2>
      <div class="header-buttons">
        <router-link to="/issued-tickets" class="back-button">
          Back to Tickets
        </router-link>
        <button @click="loadTicket" :disabled="isLoading" class="refresh-button">
          {{ isLoading ? 'Loading...' : 'Refresh' }}
        </button>
      </div>
    </div>

    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>

    <div v-if="isLoading && !ticket" class="loading-message">
      Loading ticket details...
    </div>

    <div v-if="ticket" class="ticket-details-content">
      <div class="ticket-section">
        <h3>Basic Information</h3>
        <div class="info-grid">
          <div class="info-item">
            <label>Event ID</label>
            <div class="info-value">{{ ticket.eventID || 'N/A' }}</div>
          </div>
          <div class="info-item">
            <label>Ticket ID</label>
            <div class="info-value ticket-id-value">{{ getTicketID() }}</div>
          </div>
          <div class="info-item">
            <label>Ticket Type</label>
            <div class="info-value">{{ getTicketType() }}</div>
          </div>
          <div class="info-item">
            <label>Valid From</label>
            <div class="info-value">{{ formatDateTime(getValidFrom()) }}</div>
          </div>
          <div class="info-item">
            <label>Valid Until</label>
            <div class="info-value">{{ formatDateTime(getValidUntil()) }}</div>
          </div>
        </div>
      </div>

      <div class="ticket-section">
        <h3>QR Code</h3>
        <div class="qr-code-display">
          <img v-if="ticket.qrCode" :src="'data:image/png;base64,' + ticket.qrCode" alt="Ticket QR Code" class="qr-code-image" />
        </div>
      </div>

      <div class="ticket-section">
        <h3>Ticket Data</h3>
        <pre class="json-output">{{ formatTicketData() }}</pre>
      </div>

      <div class="ticket-section">
        <h3>Signature</h3>
        <div class="signature-display">
          <textarea
            :value="ticket.signature || 'N/A'"
            readonly
            class="signature-textarea"
            rows="3"
          ></textarea>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { getTicketByID } from '../services/ticketApi'

export default {
  name: 'TicketDetails',
  data() {
    return {
      ticket: null,
      isLoading: false,
      errorMessage: '',
      ticketID: ''
    }
  },
  mounted() {
    this.loadTicketIDFromRoute()
    this.loadTicket()
  },
  watch: {
    '$route'(to) {
      this.loadTicketIDFromRoute()
      this.loadTicket()
    }
  },
  methods: {
    loadTicketIDFromRoute() {
      this.ticketID = this.$route.query.ticketID || ''
    },
    async loadTicket() {
      if (!this.ticketID) {
        this.errorMessage = 'Ticket ID is required'
        return
      }

      this.isLoading = true
      this.errorMessage = ''
      this.ticket = null

      try {
        const response = await getTicketByID(this.ticketID)
        this.ticket = response.ticket
        if (!this.ticket) {
          this.errorMessage = 'Ticket not found'
        }
      } catch (error) {
        this.errorMessage = `Failed to load ticket: ${error.message}`
        this.ticket = null
      } finally {
        this.isLoading = false
      }
    },
    getTicketID() {
      if (!this.ticket || !this.ticket.data) return 'N/A'
      return this.ticket.data.ticketId || this.ticket.data.ticketID || 'N/A'
    },
    getTicketType() {
      if (!this.ticket || !this.ticket.data) return 'N/A'
      return this.ticket.data.ticketType || 'N/A'
    },
    getValidFrom() {
      if (!this.ticket || !this.ticket.data) return ''
      return this.ticket.data.validFrom || ''
    },
    getValidUntil() {
      if (!this.ticket || !this.ticket.data) return ''
      return this.ticket.data.validUntil || ''
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
    formatTicketData() {
      if (!this.ticket || !this.ticket.data) return 'N/A'
      return JSON.stringify(this.ticket.data, null, 2)
    }
  }
}
</script>

<style scoped>
.ticket-details-container {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.ticket-details-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.ticket-details-header h2 {
  color: #2c3e50;
  font-size: 1.5rem;
  margin: 0;
}

.header-buttons {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}

.back-button {
  padding: 0.5rem 1rem;
  background-color: #95a5a6;
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

.back-button:hover {
  background-color: #7f8c8d;
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

.ticket-details-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.ticket-section {
  border-bottom: 1px solid #eee;
  padding-bottom: 1.5rem;
}

.ticket-section:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.ticket-section h3 {
  color: #2c3e50;
  font-size: 1.2rem;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #3498db;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-item label {
  font-weight: 600;
  color: #666;
  font-size: 0.9rem;
}

.info-value {
  color: #333;
  font-size: 1rem;
}

.ticket-id-value {
  font-family: 'Courier New', monospace;
  font-size: 0.95rem;
  word-break: break-all;
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
  max-width: 300px;
  max-height: 300px;
  display: block;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 1rem;
  background-color: white;
}

.json-output {
  background-color: #f5f5f5;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 1rem;
  font-family: 'Courier New', monospace;
  font-size: 0.9rem;
  overflow-x: auto;
  white-space: pre-wrap;
  word-wrap: break-word;
  max-height: 400px;
  overflow-y: auto;
  margin: 0;
}

.signature-display {
  margin-top: 0.5rem;
}

.signature-textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.9rem;
  font-family: 'Courier New', monospace;
  background-color: #f9f9f9;
  resize: vertical;
  box-sizing: border-box;
  color: #333;
}
</style>

