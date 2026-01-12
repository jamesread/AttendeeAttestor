<template>
  <div class="event-form-container">
    <form @submit.prevent="handleSubmit" class="event-form">
      <h2>Create New Event</h2>
      
      <div class="form-group">
        <label for="eventID">Event ID</label>
        <input
          id="eventID"
          name="eventID"
          type="text"
          v-model="formData.eventID"
          required
          placeholder="Enter unique event identifier"
        />
      </div>

      <div class="form-group">
        <label for="validFrom">Valid From</label>
        <input
          id="validFrom"
          name="validFrom"
          type="datetime-local"
          v-model="formData.validFrom"
        />
      </div>

      <div class="form-group">
        <label for="validUntil">Valid Until</label>
        <input
          id="validUntil"
          name="validUntil"
          type="datetime-local"
          v-model="formData.validUntil"
        />
      </div>

      <button type="submit" :disabled="isSubmitting" class="submit-button">
        {{ isSubmitting ? 'Creating...' : 'Create Event' }}
      </button>
    </form>

    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>

    <div v-if="createdEvent" class="event-result">
      <h3>Event Created Successfully</h3>
      <div class="event-details">
        <p><strong>Event ID:</strong> {{ createdEvent.eventID }}</p>
        <p><strong>Public Key:</strong> {{ createdEvent.publicKey }}</p>
      </div>
      <div class="qr-code-display">
        <h4>Scanner Configuration QR Code</h4>
        <p class="qr-instructions">Scan this QR code with scanner apps to configure them for this event:</p>
        <img :src="'data:image/png;base64,' + createdEvent.qrCode" alt="Scanner Setup QR Code" class="qr-code-image" />
      </div>
    </div>
  </div>
</template>

<script>
import { createEvent } from '../services/eventApi'

export default {
  name: 'CreateEventForm',
  data() {
    return {
      formData: {
        eventID: '',
        validFrom: '',
        validUntil: ''
      },
      isSubmitting: false,
      errorMessage: '',
      createdEvent: null
    }
  },
  methods: {
    async handleSubmit() {
      this.isSubmitting = true
      this.errorMessage = ''
      this.createdEvent = null

      try {
        const response = await createEvent(
          this.formData.eventID,
          this.formData.validFrom,
          this.formData.validUntil
        )
        this.createdEvent = response
        this.resetForm()
      } catch (error) {
        this.errorMessage = `Failed to create event: ${error.message}`
      } finally {
        this.isSubmitting = false
      }
    },
    resetForm() {
      this.formData = {
        eventID: '',
        validFrom: '',
        validUntil: ''
      }
    }
  }
}
</script>

<style scoped>
.event-form-container {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
}

.event-form h2 {
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

.submit-button {
  width: 100%;
  padding: 0.75rem;
  background-color: #27ae60;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.submit-button:hover:not(:disabled) {
  background-color: #229954;
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

.event-result {
  margin-top: 2rem;
  padding: 1.5rem;
  background-color: #f9f9f9;
  border-radius: 4px;
}

.event-result h3 {
  margin-bottom: 1rem;
  color: #27ae60;
}

.event-details {
  margin-bottom: 1.5rem;
}

.event-details p {
  margin: 0.5rem 0;
  color: #333;
}

.qr-code-display {
  margin-top: 1.5rem;
  padding: 1rem;
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.qr-code-display h4 {
  margin-bottom: 0.5rem;
  color: #2c3e50;
}

.qr-instructions {
  margin-bottom: 1rem;
  color: #666;
  font-size: 0.9rem;
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
</style>

