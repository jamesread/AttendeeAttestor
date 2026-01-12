<template>
  <div class="diagnostic-container">
    <h2>Diagnostic Tool</h2>
    
    <div class="tabs">
      <button 
        @click="activeTab = 'qr'" 
        :class="['tab-button', { active: activeTab === 'qr' }]"
      >
        QR Code Upload
      </button>
      <button 
        @click="activeTab = 'cbor'" 
        :class="['tab-button', { active: activeTab === 'cbor' }]"
      >
        CBOR Decoder
      </button>
      <button 
        @click="activeTab = 'base64'" 
        :class="['tab-button', { active: activeTab === 'base64' }]"
      >
        Base64 Decoder
      </button>
    </div>

    <div v-if="activeTab === 'qr'" class="tab-content">
      <p class="description">Upload a QR code image to extract and display the ticket ID</p>
      
      <div class="input-section">
        <label for="qr-file-input">QR Code Image</label>
        <input
          id="qr-file-input"
          type="file"
          accept="image/*"
          @change="handleQRUpload"
          class="file-input"
        />
        <div v-if="qrPreview" class="qr-preview-section">
          <label>Preview</label>
          <img :src="qrPreview" alt="QR Code Preview" class="qr-preview-image" />
        </div>
      </div>

      <div v-if="qrErrorMessage" class="error-message">
        {{ qrErrorMessage }}
      </div>

      <div v-if="qrRawContent" class="validation-status-section">
        <h3>Validation Status</h3>
        <div class="validation-steps">
          <div class="validation-step">
            <span class="validation-icon" :class="{ 'valid': qrCborValid === true, 'invalid': qrCborValid === false }">
              <span v-if="qrCborValid === true">✓</span>
              <span v-else-if="qrCborValid === false">✗</span>
              <span v-else>○</span>
            </span>
            <span class="validation-label">CBOR validation</span>
          </div>
          <div class="validation-step">
            <span class="validation-icon" :class="{ 'valid': qrCoseValid === true, 'invalid': qrCoseValid === false }">
              <span v-if="qrCoseValid === true">✓</span>
              <span v-else-if="qrCoseValid === false">✗</span>
              <span v-else>○</span>
            </span>
            <span class="validation-label">COSE validation</span>
          </div>
          <div class="validation-step">
            <span class="validation-icon" :class="{ 'valid': qrTicketIdFound === true, 'invalid': qrTicketIdFound === false }">
              <span v-if="qrTicketIdFound === true">✓</span>
              <span v-else-if="qrTicketIdFound === false">✗</span>
              <span v-else>○</span>
            </span>
            <span class="validation-label">Ticket ID found</span>
          </div>
        </div>
      </div>

      <div v-if="ticketID" class="ticket-id-section">
        <h3>Ticket ID</h3>
        <router-link :to="`/ticket-details?ticketID=${encodeURIComponent(ticketID)}`" class="ticket-id-display ticket-id-link">
          {{ ticketID }}
        </router-link>
        <button @click="copyTicketID" class="copy-button">Copy Ticket ID</button>
      </div>

      <div v-if="qrRawContent" class="raw-content-section">
        <label>Raw QR Code Content</label>
        <textarea
          v-model="qrRawContent"
          readonly
          class="raw-content-textarea"
          rows="4"
        ></textarea>
        <div class="button-group">
          <button @click="copyRawContent" class="copy-button">Copy All</button>
          <button @click="sendToCBOR" class="copy-button">To CBOR</button>
        </div>
      </div>
    </div>

    <div v-if="activeTab === 'cbor'" class="tab-content">
      <p class="description">Decode CBOR (Concise Binary Object Representation) data and view the decoded structure as JSON. Paste base64url-encoded or raw CBOR bytes.</p>
      
      <div class="input-section">
        <label for="cbor-input">CBOR Data</label>
        <textarea
          id="cbor-input"
          v-model="cborInput"
          class="cbor-textarea"
          placeholder="Paste base64url-encoded CBOR data here (e.g., from QR code) or raw CBOR bytes"
          rows="6"
        ></textarea>
      </div>

      <div class="button-section">
        <button @click="decodeCBOR" :disabled="!cborInput || isDecoding" class="decode-button">
          {{ isDecoding ? 'Decoding...' : 'Decode CBOR' }}
        </button>
        <button @click="clearCBORInput" class="clear-button">Clear</button>
      </div>

      <div v-if="cborErrorMessage" class="error-message">
        {{ cborErrorMessage }}
      </div>

      <div v-if="isCOSESign1" class="output-section">
        <h3>COSE_Sign1 Structure Detected</h3>
        
        <div class="cose-section">
          <label for="cose-protected">Protected Headers (CBOR bytes, base64)</label>
          <textarea
            id="cose-protected"
            :value="coseProtectedBase64"
            readonly
            class="cose-textarea"
            rows="2"
          ></textarea>
        </div>

        <div v-if="coseProtectedDecoded" class="cose-section">
          <label>Protected Headers (Decoded)</label>
          <pre class="json-output-small">{{ coseProtectedDecoded }}</pre>
        </div>

        <div class="cose-section">
          <label for="cose-unprotected">Unprotected Headers</label>
          <textarea
            id="cose-unprotected"
            :value="coseUnprotected"
            readonly
            class="cose-textarea"
            rows="3"
          ></textarea>
        </div>

        <div class="cose-section">
          <label for="cose-payload">Payload (CBOR bytes, base64)</label>
          <textarea
            id="cose-payload"
            :value="cosePayloadBase64"
            readonly
            class="cose-textarea"
            rows="2"
          ></textarea>
        </div>

        <div v-if="cosePayloadDecoded" class="cose-section">
          <label>Payload (Decoded)</label>
          <pre class="json-output-small">{{ cosePayloadDecoded }}</pre>
        </div>

        <div v-if="coseTicketId" class="cose-section">
          <label for="cose-ticket-id">Ticket ID</label>
          <input
            id="cose-ticket-id"
            type="text"
            :value="coseTicketId"
            readonly
            class="cose-input"
          />
        </div>

        <div v-if="coseTicketIdError" class="error-message">
          {{ coseTicketIdError }}
        </div>

        <div class="cose-section">
          <label for="cose-signature">Signature (base64)</label>
          <input
            id="cose-signature"
            type="text"
            :value="coseSignature"
            readonly
            class="cose-input"
          />
        </div>
      </div>

      <div v-if="decodedJSON && !isCOSESign1" class="output-section">
        <label>Decoded CBOR Data (JSON representation)</label>
        <pre class="json-output">{{ decodedJSON }}</pre>
        <button @click="copyJSON" class="copy-button">Copy JSON</button>
      </div>
    </div>

    <div v-if="activeTab === 'base64'" class="tab-content">
      <p class="description">Paste base64 or base64url-encoded data to decode and view as text or hex</p>
      
      <div class="input-section">
        <label for="base64-input">Base64 Data</label>
        <textarea
          id="base64-input"
          v-model="base64Input"
          class="cbor-textarea"
          placeholder="Paste base64 or base64url-encoded data here"
          rows="6"
        ></textarea>
      </div>

      <div class="button-section">
        <button @click="decodeBase64" :disabled="!base64Input || isDecodingBase64" class="decode-button">
          {{ isDecodingBase64 ? 'Decoding...' : 'Decode Base64' }}
        </button>
        <button @click="clearBase64Input" class="clear-button">Clear</button>
      </div>

      <div v-if="base64ErrorMessage" class="error-message">
        {{ base64ErrorMessage }}
      </div>

      <div v-if="decodedBase64" class="output-section">
        <label>Decoded Output</label>
        <div class="output-format-selector">
          <label>
            <input type="radio" v-model="base64OutputFormat" value="text" />
            Text (UTF-8)
          </label>
          <label>
            <input type="radio" v-model="base64OutputFormat" value="string" />
            String (with escapes)
          </label>
          <label>
            <input type="radio" v-model="base64OutputFormat" value="hex" />
            Hex
          </label>
          <label>
            <input type="radio" v-model="base64OutputFormat" value="bytes" />
            Byte Array
          </label>
        </div>
        <pre class="json-output">{{ decodedBase64 }}</pre>
        <button @click="copyBase64" class="copy-button">Copy Output</button>
      </div>
    </div>
  </div>
</template>

<script>
import jsQR from 'jsqr'

export default {
  name: 'CBORDiagnostic',
  data() {
    return {
      activeTab: 'qr',
      cborInput: '',
      decodedJSON: '',
      cborErrorMessage: '',
      isDecoding: false,
      isCOSESign1: false,
      coseProtectedBase64: '',
      coseProtectedDecoded: '',
      coseUnprotected: '',
      cosePayloadBase64: '',
      cosePayloadDecoded: '',
      coseSignature: '',
      coseTicketId: '',
      coseTicketIdError: '',
      base64Input: '',
      decodedBase64: '',
      base64ErrorMessage: '',
      isDecodingBase64: false,
      base64OutputFormat: 'text',
      qrPreview: null,
      qrErrorMessage: '',
      qrRawContent: '',
      ticketID: '',
      qrTicketIdError: '',
      qrDecodedJSON: '',
      isProcessingQR: false,
      qrCborValid: null,
      qrCoseValid: null,
      qrTicketIdFound: null
    }
  },
  methods: {
    async decodeCBOR() {
      this.isDecoding = true
      this.cborErrorMessage = ''
      this.decodedJSON = ''
      this.isCOSESign1 = false
      this.coseProtectedBase64 = ''
      this.coseProtectedDecoded = ''
      this.coseUnprotected = ''
      this.cosePayloadBase64 = ''
      this.cosePayloadDecoded = ''
      this.coseSignature = ''
      this.coseTicketId = ''
      this.coseTicketIdError = ''

      try {
        if (this.cborInput.trim() === '') {
          throw new Error('Please enter CBOR data')
        }

        const response = await fetch('http://localhost:8080/decode-cbor', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ data: this.cborInput.trim() })
        })

        if (!response.ok) {
          const errorResult = await response.json().catch(() => ({ error: `HTTP ${response.status}: ${response.statusText}` }))
          throw new Error(errorResult.error || `HTTP ${response.status}: ${response.statusText}`)
        }

        const result = await response.json()

        if (result.error) {
          throw new Error(result.error)
        }

        if (!result.decoded) {
          throw new Error('No decoded data in response')
        }

        this.decodedJSON = JSON.stringify(result.decoded, null, 2)

        if (this.isCOSESign1Structure(result.decoded)) {
          await this.parseCOSESign1(result.decoded)
        }
      } catch (error) {
        this.cborErrorMessage = `Decode error: ${error.message}`
      } finally {
        this.isDecoding = false
      }
    },
    isCOSESign1Structure(decoded) {
      if (!Array.isArray(decoded) || decoded.length !== 4) {
        return false
      }

      const [protectedHeaders, unprotectedHeaders, payload, signature] = decoded

      const isProtectedValid = typeof protectedHeaders === 'string' || Array.isArray(protectedHeaders)
      const isPayloadValid = typeof payload === 'string' || Array.isArray(payload)
      const isSignatureValid = typeof signature === 'string' || Array.isArray(signature)

      return isProtectedValid && isPayloadValid && isSignatureValid
    },
    async parseCOSESign1(decoded) {
      this.isCOSESign1 = true

      const [protectedHeaders, unprotectedHeaders, payload, signature] = decoded

      const protectedBase64 = this.valueToBase64(protectedHeaders)
      this.coseProtectedBase64 = protectedBase64

      const unprotectedJSON = JSON.stringify(unprotectedHeaders, null, 2)
      this.coseUnprotected = unprotectedJSON || '{}'

      const payloadBase64 = this.valueToBase64(payload)
      this.cosePayloadBase64 = payloadBase64

      const signatureBase64 = this.valueToBase64(signature)
      this.coseSignature = signatureBase64

      if (protectedBase64) {
        try {
          const protectedDecoded = await this.decodeNestedCBOR(protectedBase64)
          if (protectedDecoded) {
            this.coseProtectedDecoded = JSON.stringify(protectedDecoded, null, 2)
          }
        } catch (error) {
          this.coseProtectedDecoded = `Error decoding protected headers: ${error.message}`
        }
      }

      if (payloadBase64) {
        try {
          const payloadDecoded = await this.decodeNestedCBOR(payloadBase64)
          if (payloadDecoded) {
            this.cosePayloadDecoded = JSON.stringify(payloadDecoded, null, 2)
            this.extractTicketId(payloadDecoded)
          }
        } catch (error) {
          this.cosePayloadDecoded = `Error decoding payload: ${error.message}`
        }
      }
    },
    extractTicketId(payloadDecoded) {
      this.coseTicketId = ''
      this.coseTicketIdError = ''

      if (!payloadDecoded || typeof payloadDecoded !== 'object') {
        this.coseTicketIdError = 'Payload is not an object'
        return
      }

      if (payloadDecoded.ticketId) {
        this.coseTicketId = String(payloadDecoded.ticketId)
      } else if (payloadDecoded.ticketID) {
        this.coseTicketId = String(payloadDecoded.ticketID)
      } else {
        this.coseTicketIdError = 'Ticket ID not found in payload'
      }
    },
    async extractTicketIdFromCOSESign1(decoded) {
      this.ticketID = ''
      this.qrTicketIdError = ''
      this.qrTicketIdFound = null

      const [protectedHeaders, unprotectedHeaders, payload, signature] = decoded

      const payloadBase64 = this.valueToBase64(payload)

      if (!payloadBase64) {
        this.qrTicketIdError = 'Payload is not available'
        this.qrTicketIdFound = false
        return
      }

      try {
        const payloadDecoded = await this.decodeNestedCBOR(payloadBase64)
        if (payloadDecoded) {
          if (!payloadDecoded || typeof payloadDecoded !== 'object') {
            this.qrTicketIdError = 'Payload is not an object'
            this.qrTicketIdFound = false
            return
          }

          if (payloadDecoded.ticketId) {
            this.ticketID = String(payloadDecoded.ticketId)
            this.qrTicketIdError = ''
            this.qrTicketIdFound = true
          } else if (payloadDecoded.ticketID) {
            this.ticketID = String(payloadDecoded.ticketID)
            this.qrTicketIdError = ''
            this.qrTicketIdFound = true
          } else {
            this.qrTicketIdError = 'Ticket ID not found in payload'
            this.qrTicketIdFound = false
          }
        } else {
          this.qrTicketIdError = 'Failed to decode payload'
          this.qrTicketIdFound = false
        }
      } catch (error) {
        this.qrTicketIdError = `Error decoding payload: ${error.message}`
        this.qrTicketIdFound = false
      }
    },
    valueToBase64(value) {
      if (typeof value === 'string') {
        return value
      }
      if (Array.isArray(value)) {
        try {
          const bytes = new Uint8Array(value)
          let binaryString = ''
          for (let i = 0; i < bytes.length; i++) {
            binaryString += String.fromCharCode(bytes[i])
          }
          return btoa(binaryString)
        } catch (error) {
          return ''
        }
      }
      return ''
    },
    async decodeNestedCBOR(base64Data) {
      try {
        const response = await fetch('http://localhost:8080/decode-cbor', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ data: base64Data })
        })

        if (!response.ok) {
          return null
        }

        const result = await response.json()
        if (result.error || !result.decoded) {
          return null
        }

        return result.decoded
      } catch (error) {
        return null
      }
    },
    decodeBase64() {
      this.isDecodingBase64 = true
      this.base64ErrorMessage = ''
      this.decodedBase64 = ''

      try {
        if (this.base64Input.trim() === '') {
          throw new Error('Please enter base64 data')
        }

        const input = this.base64Input.trim()
        let decodedBytes

        try {
          decodedBytes = this.base64URLDecode(input)
        } catch (e1) {
          try {
            decodedBytes = this.base64Decode(input)
          } catch (e2) {
            throw new Error('Unable to decode input. Please provide valid base64 or base64url data.')
          }
        }

        if (this.base64OutputFormat === 'hex') {
          this.decodedBase64 = Array.from(decodedBytes)
            .map(b => b.toString(16).padStart(2, '0'))
            .join(' ')
        } else if (this.base64OutputFormat === 'bytes') {
          this.decodedBase64 = '[' + Array.from(decodedBytes).join(', ') + ']'
        } else if (this.base64OutputFormat === 'string') {
          const decoder = new TextDecoder('utf-8', { fatal: false })
          const text = decoder.decode(decodedBytes)
          if (this.isValidUTF8(decodedBytes)) {
            this.decodedBase64 = this.escapeString(text)
          } else {
            this.decodedBase64 = `(Binary data - ${decodedBytes.length} bytes)\n` +
              Array.from(decodedBytes)
                .map(b => b.toString(16).padStart(2, '0'))
                .join(' ')
          }
        } else {
          const decoder = new TextDecoder('utf-8', { fatal: false })
          const text = decoder.decode(decodedBytes)
          if (this.isValidUTF8(decodedBytes)) {
            this.decodedBase64 = text
          } else {
            this.decodedBase64 = `(Binary data - ${decodedBytes.length} bytes)\n` +
              Array.from(decodedBytes)
                .map(b => b.toString(16).padStart(2, '0'))
                .join(' ')
          }
        }
      } catch (error) {
        this.base64ErrorMessage = `Decode error: ${error.message}`
      } finally {
        this.isDecodingBase64 = false
      }
    },
    base64URLDecode(str) {
      let base64 = str.replace(/-/g, '+').replace(/_/g, '/')
      while (base64.length % 4) {
        base64 += '='
      }
      const binaryString = atob(base64)
      const bytes = new Uint8Array(binaryString.length)
      for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i)
      }
      return bytes
    },
    base64Decode(str) {
      const binaryString = atob(str)
      const bytes = new Uint8Array(binaryString.length)
      for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i)
      }
      return bytes
    },
    isValidUTF8(bytes) {
      try {
        const decoder = new TextDecoder('utf-8', { fatal: true })
        decoder.decode(bytes)
        return true
      } catch {
        return false
      }
    },
    escapeString(str) {
      let result = ''
      for (let i = 0; i < str.length; i++) {
        const char = str[i]
        const code = str.charCodeAt(i)
        if (code < 32 || code === 127) {
          if (code === 9) {
            result += '\\t'
          } else if (code === 10) {
            result += '\\n'
          } else if (code === 13) {
            result += '\\r'
          } else {
            result += `\\x${code.toString(16).padStart(2, '0')}`
          }
        } else if (code === 92) {
          result += '\\\\'
        } else if (code === 34) {
          result += '\\"'
        } else if (code === 39) {
          result += "\\'"
        } else {
          result += char
        }
      }
      return result
    },
    clearCBORInput() {
      this.cborInput = ''
      this.decodedJSON = ''
      this.cborErrorMessage = ''
      this.isCOSESign1 = false
      this.coseProtectedBase64 = ''
      this.coseProtectedDecoded = ''
      this.coseUnprotected = ''
      this.cosePayloadBase64 = ''
      this.cosePayloadDecoded = ''
      this.coseSignature = ''
      this.coseTicketId = ''
      this.coseTicketIdError = ''
    },
    clearBase64Input() {
      this.base64Input = ''
      this.decodedBase64 = ''
      this.base64ErrorMessage = ''
    },
    async copyJSON() {
      try {
        await navigator.clipboard.writeText(this.decodedJSON)
      } catch (err) {
        this.cborErrorMessage = 'Failed to copy to clipboard'
      }
    },
    async copyBase64() {
      try {
        await navigator.clipboard.writeText(this.decodedBase64)
      } catch (err) {
        this.base64ErrorMessage = 'Failed to copy to clipboard'
      }
    },
    async handleQRUpload(event) {
      const file = event.target.files[0]
      if (!file) {
        return
      }

      this.isProcessingQR = true
      this.qrErrorMessage = ''
      this.qrRawContent = ''
      this.ticketID = ''
      this.qrTicketIdError = ''
      this.qrDecodedJSON = ''
      this.qrPreview = null
      this.qrCborValid = null
      this.qrCoseValid = null
      this.qrTicketIdFound = null

      try {
        const imageUrl = URL.createObjectURL(file)
        this.qrPreview = imageUrl

        const img = new Image()
        img.onload = async () => {
          try {
            const canvas = document.createElement('canvas')
            canvas.width = img.width
            canvas.height = img.height
            const ctx = canvas.getContext('2d')
            ctx.drawImage(img, 0, 0)

            const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height)
            const code = jsQR(imageData.data, imageData.width, imageData.height)

            if (!code) {
              throw new Error('Could not decode QR code from image')
            }

            const qrData = code.data
            this.qrRawContent = qrData

            this.qrCborValid = null
            this.qrCoseValid = null
            this.qrTicketIdFound = null

            const response = await fetch('http://localhost:8080/decode-cbor', {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json'
              },
              body: JSON.stringify({ data: qrData })
            })

            if (!response.ok) {
              const errorResult = await response.json().catch(() => ({ error: `HTTP ${response.status}: ${response.statusText}` }))
              this.qrCborValid = false
              throw new Error(errorResult.error || `HTTP ${response.status}: ${response.statusText}`)
            }

            const result = await response.json()

            if (result.error) {
              this.qrCborValid = false
              throw new Error(result.error)
            }

            if (!result.decoded) {
              this.qrCborValid = false
              throw new Error('No decoded data in response')
            }

            this.qrCborValid = true
            this.qrDecodedJSON = JSON.stringify(result.decoded, null, 2)

            if (this.isCOSESign1Structure(result.decoded)) {
              this.qrCoseValid = true
              await this.extractTicketIdFromCOSESign1(result.decoded)
            } else {
              this.qrCoseValid = false
              if (result.decoded.ticketId) {
                this.ticketID = String(result.decoded.ticketId)
                this.qrTicketIdError = ''
                this.qrTicketIdFound = true
              } else if (result.decoded.ticketID) {
                this.ticketID = String(result.decoded.ticketID)
                this.qrTicketIdError = ''
                this.qrTicketIdFound = true
              } else {
                this.ticketID = ''
                this.qrTicketIdError = 'Ticket ID not found in decoded data'
                this.qrTicketIdFound = false
              }
            }
          } catch (error) {
            this.qrErrorMessage = `Error processing QR code: ${error.message}`
          } finally {
            this.isProcessingQR = false
          }
        }

        img.onerror = () => {
          this.qrErrorMessage = 'Failed to load image'
          this.isProcessingQR = false
        }

        img.src = imageUrl
      } catch (error) {
        this.qrErrorMessage = `Error: ${error.message}`
        this.isProcessingQR = false
      }
    },
    async copyTicketID() {
      try {
        await navigator.clipboard.writeText(this.ticketID)
      } catch (err) {
        this.qrErrorMessage = 'Failed to copy ticket ID to clipboard'
      }
    },
    async copyQRJSON() {
      try {
        await navigator.clipboard.writeText(this.qrDecodedJSON)
      } catch (err) {
        this.qrErrorMessage = 'Failed to copy JSON to clipboard'
      }
    },
    async copyRawContent() {
      try {
        await navigator.clipboard.writeText(this.qrRawContent)
      } catch (err) {
        this.qrErrorMessage = 'Failed to copy raw content to clipboard'
      }
    },
    sendToCBOR() {
      this.activeTab = 'cbor'
      this.cborInput = this.qrRawContent
      this.cborErrorMessage = ''
    }
  }
}
</script>

<style scoped>
.diagnostic-container {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.diagnostic-container h2 {
  color: #2c3e50;
  font-size: 1.5rem;
  margin-bottom: 0.5rem;
}

.tabs {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
  border-bottom: 2px solid #ddd;
}

.tab-button {
  padding: 0.75rem 1.5rem;
  background-color: transparent;
  border: none;
  border-bottom: 3px solid transparent;
  color: #666;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  margin-bottom: -2px;
}

.tab-button:hover {
  color: #2c3e50;
  background-color: #f5f5f5;
}

.tab-button.active {
  color: #3498db;
  border-bottom-color: #3498db;
}

.tab-content {
  margin-top: 1rem;
}

.description {
  color: #666;
  margin-bottom: 1.5rem;
  font-size: 0.9rem;
}

.input-section {
  margin-bottom: 1rem;
}

.input-section label {
  display: block;
  font-weight: 500;
  color: #333;
  margin-bottom: 0.5rem;
}

.cbor-textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.9rem;
  font-family: 'Courier New', monospace;
  resize: vertical;
  box-sizing: border-box;
}

.cbor-textarea:focus {
  outline: none;
  border-color: #3498db;
}

.button-section {
  display: flex;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.decode-button {
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

.decode-button:hover:not(:disabled) {
  background-color: #2980b9;
}

.decode-button:disabled {
  background-color: #95a5a6;
  cursor: not-allowed;
}

.clear-button {
  padding: 0.75rem 1.5rem;
  background-color: #95a5a6;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.clear-button:hover {
  background-color: #7f8c8d;
}

.error-message {
  padding: 1rem;
  background-color: #fee;
  border: 1px solid #fcc;
  border-radius: 4px;
  color: #c33;
  margin-bottom: 1rem;
}

.output-section {
  margin-top: 1.5rem;
}

.output-section label {
  display: block;
  font-weight: 500;
  color: #333;
  margin-bottom: 0.5rem;
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
  max-height: 500px;
  overflow-y: auto;
  margin-bottom: 0.75rem;
}

.copy-button {
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

.copy-button:hover {
  background-color: #229954;
}

.output-format-selector {
  display: flex;
  gap: 1rem;
  margin-bottom: 0.75rem;
  padding: 0.5rem;
  background-color: #f9f9f9;
  border-radius: 4px;
}

.output-format-selector label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-weight: normal;
  margin: 0;
}

.output-format-selector input[type="radio"] {
  margin: 0;
  cursor: pointer;
}

.file-input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
}

.qr-preview-section {
  margin-top: 1rem;
}

.qr-preview-section label {
  display: block;
  font-weight: 500;
  color: #333;
  margin-bottom: 0.5rem;
}

.qr-preview-image {
  max-width: 300px;
  max-height: 300px;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 0.5rem;
  background-color: white;
}

.ticket-id-section {
  margin-top: 1.5rem;
  padding: 1.5rem;
  background-color: #e8f5e9;
  border: 2px solid #4caf50;
  border-radius: 4px;
}

.ticket-id-section h3 {
  margin: 0 0 0.75rem 0;
  color: #2c3e50;
  font-size: 1.1rem;
}

.ticket-id-display {
  font-family: 'Courier New', monospace;
  font-size: 1.2rem;
  font-weight: 600;
  color: #2c3e50;
  padding: 0.75rem;
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin-bottom: 0.75rem;
  word-break: break-all;
  display: block;
  text-decoration: none;
  transition: all 0.2s;
}

.ticket-id-link {
  color: #3498db;
  cursor: pointer;
}

.ticket-id-link:hover {
  color: #2980b9;
  border-color: #3498db;
  background-color: #f0f8ff;
}

.raw-content-section {
  margin-top: 1.5rem;
}

.raw-content-section label {
  display: block;
  font-weight: 500;
  color: #333;
  margin-bottom: 0.5rem;
}

.raw-content-textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.9rem;
  font-family: 'Courier New', monospace;
  background-color: #f5f5f5;
  resize: vertical;
  box-sizing: border-box;
  margin-bottom: 0.75rem;
}

.validation-status-section {
  margin-top: 1.5rem;
  padding: 1.5rem;
  background-color: #f9f9f9;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.validation-status-section h3 {
  margin: 0 0 1rem 0;
  color: #2c3e50;
  font-size: 1.1rem;
}

.validation-steps {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.validation-step {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.validation-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  font-size: 1rem;
  font-weight: bold;
  flex-shrink: 0;
}

.validation-icon.valid {
  background-color: #4caf50;
  color: white;
}

.validation-icon.invalid {
  background-color: #f44336;
  color: white;
}

.validation-icon:not(.valid):not(.invalid) {
  background-color: #e0e0e0;
  color: #999;
}

.validation-label {
  color: #333;
  font-size: 0.95rem;
}

.button-group {
  display: flex;
  gap: 0.75rem;
}

.output-section h3 {
  color: #2c3e50;
  font-size: 1.2rem;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #3498db;
}

.cose-section {
  margin-bottom: 1.5rem;
}

.cose-section label {
  display: block;
  font-weight: 500;
  color: #333;
  margin-bottom: 0.5rem;
  font-size: 0.95rem;
}

.cose-textarea {
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

.cose-input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.9rem;
  font-family: 'Courier New', monospace;
  background-color: #f9f9f9;
  box-sizing: border-box;
  color: #333;
}

.json-output-small {
  background-color: #f5f5f5;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 0.75rem;
  font-family: 'Courier New', monospace;
  font-size: 0.85rem;
  overflow-x: auto;
  white-space: pre-wrap;
  word-wrap: break-word;
  max-height: 300px;
  overflow-y: auto;
  margin-top: 0.5rem;
}
</style>

