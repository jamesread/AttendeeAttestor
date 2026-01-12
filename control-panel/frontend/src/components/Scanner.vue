<template>
  <div class="scanner-container">
    <div v-if="!isConfigured" class="setup-section">
      <h3>Scanner Setup</h3>
      <p class="description">Scan the scanner setup QR code to configure this scanner with the event's public key.</p>
      
      <div v-if="setupError" class="error-message">
        {{ setupError }}
      </div>

      <div class="camera-section">
        <video
          ref="setupVideo"
          autoplay
          playsinline
          class="camera-video"
        ></video>
        <canvas ref="setupCanvas" class="hidden-canvas"></canvas>
      </div>

      <div class="button-group">
        <button @click="startSetupScan" :disabled="isScanningSetup" class="primary-button">
          {{ isScanningSetup ? 'Scanning...' : 'Start Setup Scan' }}
        </button>
        <button @click="stopSetupScan" v-if="isScanningSetup" class="secondary-button">
          Stop
        </button>
      </div>

      <div v-if="scannerConfig" class="config-display">
        <h4>Configuration Loaded</h4>
        <p><strong>Event ID:</strong> {{ scannerConfig.eventID }}</p>
        <p><strong>Public Key:</strong> {{ scannerConfig.publicKey.substring(0, 20) }}...</p>
        <button @click="clearConfiguration" class="secondary-button">Clear Configuration</button>
      </div>
    </div>

    <div v-else class="ticket-scan-section">
      <div class="status-header">
        <h3 @click="toggleEventMenu" class="event-title-clickable">Event: {{ scannerConfig.eventID }}</h3>
      </div>

      <div v-if="showChangeEventDialog" class="dialog-overlay" @click.self="closeEventDialog">
        <div class="dialog-content">
          <h3>Change Event</h3>
          <p>Are you sure you want to change the event? This will clear the current scanner configuration.</p>
          <div class="dialog-buttons">
            <button @click="clearConfiguration" class="primary-button">
              Change Event
            </button>
            <button @click="closeEventDialog" class="secondary-button">
              Cancel
            </button>
          </div>
        </div>
      </div>

      <div class="camera-section">
        <video
          ref="ticketVideo"
          autoplay
          playsinline
          class="camera-video"
          :class="{ 'camera-hidden': lastResult && (lastResult.valid || !lastResult.valid) }"
        ></video>
        <canvas ref="ticketCanvas" class="hidden-canvas"></canvas>
        
        <div 
          v-if="lastResult && lastResult.valid" 
          @click="dismissSuccess" 
          class="success-overlay"
        >
          <div class="success-content">
            <div class="success-icon">✓</div>
            <h3>VALID TICKET</h3>
            <div v-if="lastResult.ticketData" class="ticket-info-overlay">
              <p><strong>Ticket ID:</strong> {{ lastResult.ticketData.ticketID || lastResult.ticketData.ticketId || 'N/A' }}</p>
              <p><strong>Ticket Type:</strong> {{ lastResult.ticketData.ticketType || 'N/A' }}</p>
              <p v-if="lastResult.ticketData.validFrom"><strong>Valid From:</strong> {{ lastResult.ticketData.validFrom }}</p>
              <p v-if="lastResult.ticketData.validUntil"><strong>Valid Until:</strong> {{ lastResult.ticketData.validUntil }}</p>
            </div>
            <p class="success-hint">Tap to continue scanning</p>
          </div>
        </div>
        
        <div 
          v-if="lastResult && !lastResult.valid" 
          @click="dismissError" 
          class="error-overlay"
        >
          <div class="error-content">
            <div class="error-icon">✗</div>
            <h3>INVALID TICKET</h3>
            <p class="error-text">{{ lastResult.error || 'Unknown error' }}</p>
            <p class="error-hint">Tap to dismiss</p>
          </div>
        </div>
      </div>

      <div class="button-group">
        <button @click="startTicketScan" v-if="!isScanningTicket" class="primary-button">
          Start Ticket Scan
        </button>
        <button @click="stopTicketScan" v-if="isScanningTicket" class="secondary-button">
          Stop
        </button>
      </div>

      <div v-if="usedTickets.length > 0" class="used-tickets-section">
        <h4>Used Tickets ({{ usedTickets.length }})</h4>
        <button @click="clearUsedTickets" class="secondary-button">Clear History</button>
        <ul class="used-tickets-list">
          <li v-for="ticketId in usedTickets.slice(-10)" :key="ticketId">{{ ticketId }}</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
import jsQR from 'jsqr'
import { verify } from '@noble/ed25519'
import CBOR from 'cbor-js'

const STORAGE_KEY_CONFIG = 'scannerConfig'
const STORAGE_KEY_USED_TICKETS = 'usedTickets'

export default {
  name: 'Scanner',
  data() {
    return {
      isConfigured: false,
      scannerConfig: null,
      isScanningSetup: false,
      isScanningTicket: false,
      showChangeEventDialog: false,
      setupError: '',
      lastResult: null,
      usedTickets: [],
      setupStream: null,
      ticketStream: null,
      setupScanInterval: null,
      ticketScanInterval: null
    }
  },
  mounted() {
    this.loadConfiguration()
    this.loadUsedTickets()
  },
  beforeUnmount() {
    this.stopSetupScan()
    this.stopTicketScan()
  },
  methods: {
    loadConfiguration() {
      const stored = localStorage.getItem(STORAGE_KEY_CONFIG)
      if (stored) {
        try {
          this.scannerConfig = JSON.parse(stored)
          this.isConfigured = true
        } catch (error) {
          console.error('Failed to load configuration:', error)
        }
      }
    },
    saveConfiguration() {
      if (this.scannerConfig) {
        localStorage.setItem(STORAGE_KEY_CONFIG, JSON.stringify(this.scannerConfig))
      }
    },
    clearConfiguration() {
      this.scannerConfig = null
      this.isConfigured = false
      this.showChangeEventDialog = false
      localStorage.removeItem(STORAGE_KEY_CONFIG)
      this.stopSetupScan()
      this.stopTicketScan()
    },
    toggleEventMenu() {
      this.showChangeEventDialog = true
    },
    closeEventDialog() {
      this.showChangeEventDialog = false
    },
    loadUsedTickets() {
      const stored = localStorage.getItem(STORAGE_KEY_USED_TICKETS)
      if (stored) {
        try {
          this.usedTickets = JSON.parse(stored)
        } catch (error) {
          console.error('Failed to load used tickets:', error)
          this.usedTickets = []
        }
      }
    },
    saveUsedTickets() {
      localStorage.setItem(STORAGE_KEY_USED_TICKETS, JSON.stringify(this.usedTickets))
    },
    addUsedTicket(ticketId) {
      if (ticketId && !this.usedTickets.includes(ticketId)) {
        this.usedTickets.push(ticketId)
        this.saveUsedTickets()
      }
    },
    clearUsedTickets() {
      this.usedTickets = []
      this.saveUsedTickets()
    },
    async startSetupScan() {
      this.setupError = ''
      
      if (!navigator.mediaDevices) {
        const isSecure = window.location.protocol === 'https:' || window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
        if (!isSecure) {
          this.setupError = 'Camera access requires HTTPS. Please access this page over HTTPS or use localhost.'
        } else {
          this.setupError = 'Camera access is not available in this browser. Please use a modern browser that supports camera access.'
        }
        return
      }

      if (!navigator.mediaDevices.getUserMedia) {
        this.setupError = 'Camera API not supported. Please use a modern browser.'
        return
      }

      try {
        const constraints = {
          video: { facingMode: 'environment' }
        }
        this.setupStream = await navigator.mediaDevices.getUserMedia(constraints)
        if (this.$refs.setupVideo) {
          this.$refs.setupVideo.srcObject = this.setupStream
          this.isScanningSetup = true
          this.setupScanInterval = setInterval(() => this.scanSetupQR(), 200)
        }
      } catch (error) {
        let errorMessage = `Failed to access camera: ${error.message}`
        if (error.name === 'NotAllowedError') {
          errorMessage = 'Camera permission denied. Please allow camera access in your browser settings.'
        } else if (error.name === 'NotFoundError') {
          errorMessage = 'No camera found. Please ensure your device has a camera.'
        } else if (error.name === 'NotReadableError') {
          errorMessage = 'Camera is already in use by another application.'
        }
        this.setupError = errorMessage
        console.error('Camera access error:', error)
      }
    },
    stopSetupScan() {
      if (this.setupStream) {
        this.setupStream.getTracks().forEach(track => track.stop())
        this.setupStream = null
      }
      if (this.setupScanInterval) {
        clearInterval(this.setupScanInterval)
        this.setupScanInterval = null
      }
      this.isScanningSetup = false
    },
    scanSetupQR() {
      if (!this.$refs.setupVideo || !this.$refs.setupCanvas) return

      const video = this.$refs.setupVideo
      const canvas = this.$refs.setupCanvas
      const context = canvas.getContext('2d')

      if (video.readyState === video.HAVE_ENOUGH_DATA) {
        canvas.width = video.videoWidth
        canvas.height = video.videoHeight
        context.drawImage(video, 0, 0, canvas.width, canvas.height)
        const imageData = context.getImageData(0, 0, canvas.width, canvas.height)
        const code = jsQR(imageData.data, imageData.width, imageData.height)

        if (code) {
          this.processSetupQR(code.data)
        }
      }
    },
    processSetupQR(data) {
      try {
        const config = JSON.parse(data)
        if (config.eventID && config.publicKey) {
          this.scannerConfig = {
            eventID: config.eventID,
            publicKey: config.publicKey
          }
          this.saveConfiguration()
          this.isConfigured = true
          this.stopSetupScan()
          this.setupError = ''
        } else {
          this.setupError = 'Invalid scanner configuration format'
        }
      } catch (error) {
        this.setupError = `Failed to parse setup QR: ${error.message}`
      }
    },
    async startTicketScan() {
      this.lastResult = null
      
      if (!navigator.mediaDevices) {
        const isSecure = window.location.protocol === 'https:' || window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
        let errorMessage = 'Camera access is not available in this browser.'
        if (!isSecure) {
          errorMessage = 'Camera access requires HTTPS. Please access this page over HTTPS or use localhost.'
        }
        this.lastResult = {
          valid: false,
          error: errorMessage
        }
        return
      }

      if (!navigator.mediaDevices.getUserMedia) {
        this.lastResult = {
          valid: false,
          error: 'Camera API not supported. Please use a modern browser.'
        }
        return
      }

      try {
        const constraints = {
          video: { facingMode: 'environment' }
        }
        this.ticketStream = await navigator.mediaDevices.getUserMedia(constraints)
        if (this.$refs.ticketVideo) {
          this.$refs.ticketVideo.srcObject = this.ticketStream
          this.isScanningTicket = true
          this.ticketScanInterval = setInterval(() => this.scanTicketQR(), 200)
        }
      } catch (error) {
        let errorMessage = `Failed to access camera: ${error.message}`
        if (error.name === 'NotAllowedError') {
          errorMessage = 'Camera permission denied. Please allow camera access in your browser settings.'
        } else if (error.name === 'NotFoundError') {
          errorMessage = 'No camera found. Please ensure your device has a camera.'
        } else if (error.name === 'NotReadableError') {
          errorMessage = 'Camera is already in use by another application.'
        }
        this.lastResult = {
          valid: false,
          error: errorMessage
        }
        console.error('Camera access error:', error)
      }
    },
    stopTicketScan() {
      if (this.ticketStream) {
        this.ticketStream.getTracks().forEach(track => track.stop())
        this.ticketStream = null
      }
      if (this.ticketScanInterval) {
        clearInterval(this.ticketScanInterval)
        this.ticketScanInterval = null
      }
      this.isScanningTicket = false
    },
    scanTicketQR() {
      if (!this.$refs.ticketVideo || !this.$refs.ticketCanvas) return

      const video = this.$refs.ticketVideo
      const canvas = this.$refs.ticketCanvas
      const context = canvas.getContext('2d')

      if (video.readyState === video.HAVE_ENOUGH_DATA) {
        canvas.width = video.videoWidth
        canvas.height = video.videoHeight
        context.drawImage(video, 0, 0, canvas.width, canvas.height)
        const imageData = context.getImageData(0, 0, canvas.width, canvas.height)
        const code = jsQR(imageData.data, imageData.width, imageData.height)

        if (code) {
          this.processTicketQR(code.data)
        }
      }
    },
    async processTicketQR(data) {
      try {
        const result = await this.verifyTicket(data)
        this.lastResult = result
        
        if (result.valid && result.ticketData) {
          const ticketId = result.ticketData.ticketID || result.ticketData.ticketId
          if (ticketId) {
            this.addUsedTicket(ticketId)
          }
        } else {
          this.lastResult = result
        }
      } catch (error) {
        this.lastResult = {
          valid: false,
          error: `Verification error: ${error.message}`
        }
      }
    },
    dismissError() {
      this.lastResult = null
    },
    dismissSuccess() {
      this.lastResult = null
    },
    async verifyTicket(qrData) {
      if (!this.scannerConfig) {
        return { valid: false, error: 'Scanner not configured' }
      }

      try {
        const coseBase64 = this.decodeBase64URL(qrData)
        const coseBytes = this.base64ToBytes(coseBase64)
        const coseArray = CBOR.decode(coseBytes.buffer)

        if (!Array.isArray(coseArray) || coseArray.length !== 4) {
          return { valid: false, error: 'Invalid COSE structure' }
        }

        const [protectedHeaders, unprotectedHeaders, payload, signature] = coseArray

        let protectedBytes
        let payloadBytes
        let signatureBytes

        if (protectedHeaders instanceof Uint8Array) {
          protectedBytes = protectedHeaders
        } else if (Array.isArray(protectedHeaders)) {
          protectedBytes = new Uint8Array(protectedHeaders)
        } else {
          return { valid: false, error: 'Invalid protected headers format' }
        }

        if (payload instanceof Uint8Array) {
          payloadBytes = payload
        } else if (Array.isArray(payload)) {
          payloadBytes = new Uint8Array(payload)
        } else {
          return { valid: false, error: 'Invalid payload format' }
        }

        if (signature instanceof Uint8Array) {
          signatureBytes = signature
        } else if (Array.isArray(signature)) {
          signatureBytes = new Uint8Array(signature)
        } else {
          return { valid: false, error: 'Invalid signature format' }
        }

        const protectedDecoded = CBOR.decode(protectedBytes.buffer)
        const payloadDecoded = CBOR.decode(payloadBytes.buffer)

        const protectedMap = this.cborMapToObject(protectedDecoded)
        
        if (!this.isCOSESign1Valid(protectedMap)) {
          console.error('Protected headers:', protectedMap, 'Decoded:', protectedDecoded)
          return { valid: false, error: 'Invalid COSE algorithm or key ID' }
        }

        const keyID = this.getKeyID(protectedMap)
        const keyIDString = String(keyID)
        const eventIDString = String(this.scannerConfig.eventID)
        const publicKeyString = String(this.scannerConfig.publicKey)
        
        if (keyIDString !== eventIDString && keyIDString !== publicKeyString) {
          console.error('Key ID mismatch:', { keyID: keyIDString, eventID: eventIDString, publicKey: publicKeyString })
          return { valid: false, error: 'Key ID mismatch' }
        }

        const sigStructure = [
          'Signature1',
          protectedBytes,
          new Uint8Array(0),
          payloadBytes
        ]

        const sigStructureEncoded = CBOR.encode(sigStructure)
        const sigStructureBytes = new Uint8Array(sigStructureEncoded)
        const publicKeyBytes = this.base64ToBytes(this.scannerConfig.publicKey)

        console.log('Signature verification:', {
          signatureLength: signatureBytes.length,
          sigStructureLength: sigStructureBytes.length,
          publicKeyLength: publicKeyBytes.length,
          protectedBytesLength: protectedBytes.length,
          payloadBytesLength: payloadBytes.length
        })

        const isValid = await verify(signatureBytes, sigStructureBytes, publicKeyBytes)

        if (!isValid) {
          console.error('Signature verification failed. Details:', {
            signature: Array.from(signatureBytes.slice(0, 10)),
            sigStructure: Array.from(sigStructureBytes.slice(0, 20)),
            protectedBytes: Array.from(protectedBytes.slice(0, 10)),
            payloadBytes: Array.from(payloadBytes.slice(0, 10))
          })
          return { valid: false, error: 'Signature verification failed' }
        }

        const ticketData = this.cborMapToObject(payloadDecoded)
        const validationError = this.validateTicketData(ticketData)

        if (validationError) {
          return { valid: false, error: validationError }
        }

        const ticketId = ticketData.ticketID || ticketData.ticketId
        if (ticketId && this.usedTickets.includes(ticketId)) {
          return { valid: false, error: 'Ticket already used' }
        }

        return { valid: true, ticketData }
      } catch (error) {
        return { valid: false, error: `Verification failed: ${error.message}` }
      }
    },
    validateTicketData(ticketData) {
      if (!ticketData || typeof ticketData !== 'object') {
        return 'Invalid ticket data'
      }

      const eventID = ticketData.eventID
      if (!eventID || eventID !== this.scannerConfig.eventID) {
        return 'Event ID mismatch'
      }

      const validFrom = ticketData.validFrom
      const validUntil = ticketData.validUntil

      if (validFrom || validUntil) {
        const now = new Date()
        const fromDate = validFrom ? this.parseDateTime(validFrom) : null
        const untilDate = validUntil ? this.parseDateTime(validUntil) : null

        if (fromDate && now < fromDate) {
          return 'Ticket not yet valid'
        }

        if (untilDate && now > untilDate) {
          return 'Ticket expired'
        }
      }

      return null
    },
    parseDateTime(dateTimeString) {
      if (!dateTimeString) return null
      
      if (dateTimeString.includes('T')) {
        const parts = dateTimeString.split('T')
        if (parts.length === 2) {
          const datePart = parts[0]
          const timePart = parts[1]
          const dateTime = `${datePart}T${timePart}:00`
          return new Date(dateTime)
        }
      }
      
      return new Date(dateTimeString)
    },
    isCOSESign1Valid(protectedHeaders) {
      if (!protectedHeaders || typeof protectedHeaders !== 'object') {
        return false
      }

      const algorithm = this.getAlgorithm(protectedHeaders)
      if (algorithm !== -8 && algorithm !== '-8') {
        return false
      }

      return true
    },
    getAlgorithm(protectedHeaders) {
      if (protectedHeaders instanceof Map) {
        return protectedHeaders.get(1)
      }
      return protectedHeaders['1'] || protectedHeaders[1] || protectedHeaders['-8']
    },
    getKeyID(protectedHeaders) {
      if (protectedHeaders instanceof Map) {
        return protectedHeaders.get(4)
      }
      const keyID = protectedHeaders['4'] || protectedHeaders[4]
      if (typeof keyID === 'string' || typeof keyID === 'number') {
        return String(keyID)
      }
      return keyID
    },
    decodeBase64URL(base64url) {
      let base64 = base64url.replace(/-/g, '+').replace(/_/g, '/')
      while (base64.length % 4) {
        base64 += '='
      }
      return base64
    },
    base64ToBytes(base64) {
      const binary = atob(base64)
      const bytes = new Uint8Array(binary.length)
      for (let i = 0; i < binary.length; i++) {
        bytes[i] = binary.charCodeAt(i)
      }
      return bytes
    },
    ensureUint8Array(value) {
      if (value instanceof Uint8Array) {
        return value
      }
      if (value instanceof ArrayBuffer) {
        return new Uint8Array(value)
      }
      if (Array.isArray(value)) {
        const isByteArray = value.length === 0 || value.every(item => typeof item === 'number' && item >= 0 && item <= 255)
        if (isByteArray) {
          return new Uint8Array(value)
        }
      }
      if (typeof value === 'string') {
        return new TextEncoder().encode(value)
      }
      return new Uint8Array(value)
    },
    cborValueToBytes(value) {
      if (value instanceof Uint8Array) {
        return value
      }
      if (value instanceof ArrayBuffer) {
        return new Uint8Array(value)
      }
      if (Array.isArray(value)) {
        const isByteArray = value.length === 0 || value.every(item => typeof item === 'number' && item >= 0 && item <= 255)
        if (isByteArray) {
          return new Uint8Array(value)
        }
        const encoded = CBOR.encode(value)
        return new Uint8Array(encoded)
      }
      if (typeof value === 'string') {
        return new TextEncoder().encode(value)
      }
      const encoded = CBOR.encode(value)
      return new Uint8Array(encoded)
    },
    cborMapToObject(cborMap) {
      if (!cborMap || (typeof cborMap !== 'object' && !(cborMap instanceof Map))) {
        return cborMap
      }

      if (cborMap instanceof Map) {
        const result = {}
        for (const [key, value] of cborMap.entries()) {
          const stringKey = String(key)
          if (typeof value === 'string' || typeof value === 'number' || typeof value === 'boolean') {
            result[stringKey] = value
          } else if (value instanceof Uint8Array) {
            result[stringKey] = new TextDecoder().decode(value)
          } else if (value instanceof Map) {
            result[stringKey] = this.cborMapToObject(value)
          } else if (Array.isArray(value)) {
            result[stringKey] = value.map(item => 
              item instanceof Uint8Array ? new TextDecoder().decode(item) : item
            )
          } else {
            result[stringKey] = value
          }
        }
        return result
      }

      const result = {}
      for (const key in cborMap) {
        if (!cborMap.hasOwnProperty(key)) continue
        const value = cborMap[key]
        if (typeof value === 'string' || typeof value === 'number' || typeof value === 'boolean') {
          result[key] = value
        } else if (value instanceof Uint8Array) {
          result[key] = new TextDecoder().decode(value)
        } else if (value instanceof Map) {
          result[key] = this.cborMapToObject(value)
        } else if (Array.isArray(value)) {
          result[key] = value.map(item => 
            item instanceof Uint8Array ? new TextDecoder().decode(item) : item
          )
        } else {
          result[key] = value
        }
      }
      return result
    }
  }
}
</script>

<style scoped>
.scanner-container {
  max-width: 800px;
  margin: 0 auto;
}

h2 {
  margin-bottom: 1.5rem;
  color: #2c3e50;
}

h3 {
  margin-bottom: 1rem;
  color: #34495e;
}

h4 {
  margin-bottom: 0.5rem;
  color: #7f8c8d;
}

.description {
  margin-bottom: 1rem;
  color: #666;
}

.setup-section,
.ticket-scan-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.status-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.event-title-clickable {
  cursor: pointer;
  user-select: none;
  transition: color 0.2s;
  margin: 0;
}

.event-title-clickable:hover {
  color: #3498db;
}

.event-title-clickable {
  cursor: pointer;
  user-select: none;
  transition: color 0.2s;
  margin: 0;
}

.event-title-clickable:hover {
  color: #3498db;
}

.camera-section {
  position: relative;
  width: 100%;
  max-width: 500px;
  margin: 0 auto;
  background: #000;
  border-radius: 8px;
  overflow: hidden;
  min-height: 300px;
}

.camera-video {
  width: 100%;
  height: auto;
  display: block;
  min-height: 300px;
}

.camera-video.camera-hidden {
  opacity: 0;
  pointer-events: none;
}

.hidden-canvas {
  display: none;
}

.button-group {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.primary-button {
  padding: 0.75rem 1.5rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.primary-button:hover:not(:disabled) {
  background-color: #2980b9;
}

.primary-button:disabled {
  background-color: #95a5a6;
  cursor: not-allowed;
}

.secondary-button {
  padding: 0.75rem 1.5rem;
  background-color: #ecf0f1;
  color: #2c3e50;
  border: 1px solid #bdc3c7;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.secondary-button:hover {
  background-color: #d5dbdb;
}

.error-message {
  padding: 1rem;
  background-color: #fee;
  color: #c33;
  border: 1px solid #fcc;
  border-radius: 4px;
}

.config-display {
  padding: 1rem;
  background-color: #e8f5e9;
  border: 1px solid #c8e6c9;
  border-radius: 4px;
}

.config-display p {
  margin: 0.5rem 0;
}

.result-section {
  padding: 1.5rem;
  border-radius: 8px;
  text-align: center;
}

.result-section.valid {
  background-color: #e8f5e9;
  border: 2px solid #4caf50;
}

.result-section.invalid {
  background-color: #ffebee;
  border: 2px solid #f44336;
}

.success-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: #e8f5e9;
  border: 3px solid #4caf50;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  z-index: 10;
}

.success-overlay:hover {
  background: #c8e6c9;
  border-color: #388e3c;
}

.success-content {
  text-align: center;
  padding: 2rem;
  color: #2e7d32;
}

.success-icon {
  font-size: 5rem;
  margin-bottom: 1rem;
  color: #4caf50;
}

.success-content h3 {
  font-size: 1.5rem;
  margin-bottom: 1rem;
  color: #2e7d32;
}

.ticket-info-overlay {
  margin: 1rem 0;
  text-align: left;
  background: white;
  padding: 1rem;
  border-radius: 4px;
  color: #2c3e50;
}

.ticket-info-overlay p {
  margin: 0.5rem 0;
  font-size: 0.95rem;
}

.success-hint {
  font-size: 0.9rem;
  color: #388e3c;
  font-style: italic;
  margin-top: 1rem;
}

.error-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: #ffebee;
  border: 3px solid #f44336;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  z-index: 10;
}

.error-overlay:hover {
  background: #ffcdd2;
  border-color: #d32f2f;
}

.error-content {
  text-align: center;
  padding: 2rem;
  color: #c62828;
}

.error-icon {
  font-size: 5rem;
  margin-bottom: 1rem;
  color: #f44336;
}

.error-content h3 {
  font-size: 1.5rem;
  margin-bottom: 1rem;
  color: #c62828;
}

.error-text {
  font-size: 1.1rem;
  margin-bottom: 1rem;
  color: #b71c1c;
}

.error-hint {
  font-size: 0.9rem;
  color: #d32f2f;
  font-style: italic;
  margin-top: 1rem;
}

.result-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
}

.result-section.valid .result-icon {
  color: #4caf50;
}

.result-section.invalid .result-icon {
  color: #f44336;
}

.result-message h3 {
  margin-bottom: 0.5rem;
  font-size: 1.5rem;
}

.ticket-info {
  margin-top: 1rem;
  text-align: left;
  background: white;
  padding: 1rem;
  border-radius: 4px;
}

.ticket-info p {
  margin: 0.5rem 0;
}

.used-tickets-section {
  margin-top: 2rem;
  padding: 1rem;
  background-color: #f5f5f5;
  border-radius: 4px;
}

.used-tickets-list {
  list-style: none;
  padding: 0;
  margin: 0.5rem 0 0 0;
  max-height: 200px;
  overflow-y: auto;
}

.used-tickets-list li {
  padding: 0.25rem 0;
  font-family: monospace;
  font-size: 0.9rem;
  color: #666;
}

.dialog-overlay {
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

.dialog-content {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  max-width: 400px;
  width: 90%;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.dialog-content h3 {
  margin: 0 0 1rem 0;
  color: #2c3e50;
  font-size: 1.5rem;
}

.dialog-content p {
  margin: 0 0 1.5rem 0;
  color: #666;
  line-height: 1.5;
}

.dialog-buttons {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}
</style>
