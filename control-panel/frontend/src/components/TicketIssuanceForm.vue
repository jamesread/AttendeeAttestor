<template>
  <div class="ticket-form-container">
    <form @submit.prevent="handleSubmit" class="ticket-form">
      <h2>Issue New Ticket</h2>
      
      <div class="form-group">
        <label for="eventID">Event ID</label>
        <span v-if="eventIDFromRoute" class="event-id-display">{{ formData.eventID }}</span>
        <input
          v-else
          id="eventID"
          name="eventID"
          type="text"
          v-model="formData.eventID"
          required
          placeholder="Enter event ID"
        />
      </div>

      <div class="form-group">
        <label for="ticketType">Ticket Type</label>
        <input
          id="ticketType"
          name="ticketType"
          type="text"
          v-model="formData.ticketType"
          placeholder="e.g., General, VIP"
        />
      </div>

      <button type="submit" :disabled="isSubmitting" class="submit-button">
        {{ isSubmitting ? 'Generating...' : 'Generate Ticket' }}
      </button>
    </form>

    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>

    <div v-if="generatedTicket" class="ticket-result">
      <h3>Generated Ticket</h3>
      <div class="qr-code-display">
        <img :src="'data:image/png;base64,' + generatedTicket.qrCode" alt="Ticket QR Code" class="qr-code-image" />
      </div>
      <div class="ticket-details">
        <p><strong>Event ID:</strong> {{ generatedTicket.data.eventID }}</p>
        <p v-if="generatedTicket.data.ticketId"><strong>Ticket ID:</strong> {{ generatedTicket.data.ticketId }}</p>
        <p v-if="generatedTicket.data.ticketType"><strong>Ticket Type:</strong> {{ generatedTicket.data.ticketType }}</p>
        <p><strong>Signature:</strong> {{ generatedTicket.signature.substring(0, 20) }}...</p>
      </div>
    </div>
  </div>
</template>

<script>
import { generateTicket } from '../services/ticketApi'

export default {
  name: 'TicketIssuanceForm',
  data() {
    return {
      formData: {
        eventID: '',
        ticketType: ''
      },
      isSubmitting: false,
      errorMessage: '',
      generatedTicket: null
    }
  },
  mounted() {
    this.loadEventIDFromRoute()
  },
  watch: {
    '$route'(to) {
      this.loadEventIDFromRoute()
    }
  },
  methods: {
    loadEventIDFromRoute() {
      const eventID = this.$route.query.eventID
      if (eventID) {
        this.formData.eventID = eventID
        this.eventIDFromRoute = true
      } else {
        this.eventIDFromRoute = false
      }
    },
    async handleSubmit() {
      this.isSubmitting = true
      this.errorMessage = ''
      this.generatedTicket = null

      try {
        const ticketData = this.buildTicketData()
        const response = await generateTicket(ticketData)
        this.generatedTicket = response
        this.resetForm()
      } catch (error) {
        this.errorMessage = `Failed to generate ticket: ${error.message}`
      } finally {
        this.isSubmitting = false
      }
    },
    buildTicketData() {
      const ticketData = {}
      this.addFieldIfPresent(ticketData, 'eventID', this.formData.eventID)
      this.addFieldIfPresent(ticketData, 'ticketType', this.formData.ticketType)
      return ticketData
    },
    addFieldIfPresent(ticketData, fieldName, fieldValue) {
      if (fieldValue) {
        ticketData[fieldName] = fieldValue
      }
    },
    resetForm() {
      const eventIDFromRoute = this.$route.query.eventID || ''
      this.formData = {
        eventID: eventIDFromRoute,
        ticketType: ''
      }
      this.eventIDFromRoute = !!eventIDFromRoute
    }
  }
}
</script>

<style scoped>
.ticket-form-container {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.ticket-form h2 {
  margin-bottom: 1.5rem;
  color: #2c3e50;
  font-size: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.form-group label {
  flex: 0 0 150px;
  font-weight: 500;
  color: #333;
  text-align: right;
}

.form-group input {
  flex: 1;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  transition: border-color 0.2s;
}

.form-group input:focus {
  outline: none;
  border-color: #3498db;
}

.event-id-display {
  flex: 1;
  padding: 0.75rem;
  background-color: #f5f5f5;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  color: #333;
  font-weight: 500;
}

.submit-button {
  width: 100%;
  padding: 0.75rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.submit-button:hover:not(:disabled) {
  background-color: #2980b9;
}

.submit-button:disabled {
  background-color: #95a5a6;
  cursor: not-allowed;
}

.error-message {
  margin-top: 1rem;
  padding: 1rem;
  background-color: #fee;
  border: 1px solid #fcc;
  border-radius: 4px;
  color: #c33;
}

.ticket-result {
  margin-top: 2rem;
  padding: 1.5rem;
  background-color: #f9f9f9;
  border-radius: 4px;
}

.ticket-result h3 {
  margin-bottom: 1rem;
  color: #2c3e50;
}

.qr-code-display {
  margin: 1rem 0;
  padding: 1rem;
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  overflow-x: auto;
}

.qr-code-image {
  max-width: 100%;
  height: auto;
  display: block;
  margin: 0 auto;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 1rem;
  background-color: white;
}

.ticket-details {
  margin-top: 1rem;
}

.ticket-details p {
  margin: 0.5rem 0;
  color: #333;
}
</style>

