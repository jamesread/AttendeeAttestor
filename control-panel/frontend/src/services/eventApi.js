const ISSUER_SERVICE_URL = 'http://localhost:8080/create-event'
const LIST_EVENTS_URL = 'http://localhost:8080/list-events'

export async function createEvent(eventID, validFrom, validUntil) {
  try {
    const requestBody = { eventID }
    if (validFrom) {
      requestBody.validFrom = validFrom
    }
    if (validUntil) {
      requestBody.validUntil = validUntil
    }
    
    const response = await fetch(ISSUER_SERVICE_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestBody)
    })

    if (!response.ok) {
      throw new Error(`Failed to create event: ${response.status} ${response.statusText}`)
    }

    return await response.json()
  } catch (error) {
    if (error.message.includes('Failed to create event')) {
      throw error
    }
    throw new Error(`Network error: ${error.message}`)
  }
}

export async function listEvents() {
  try {
    const response = await fetch(LIST_EVENTS_URL, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      throw new Error(`Failed to list events: ${response.status} ${response.statusText}`)
    }

    return await response.json()
  } catch (error) {
    if (error.message.includes('Failed to list events')) {
      throw error
    }
    throw new Error(`Network error: ${error.message}`)
  }
}

export async function getScannerQR(eventID) {
  try {
    const response = await fetch(`http://localhost:8080/get-scanner-qr/${eventID}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      if (response.status === 404) {
        throw new Error(`Event not found: ${eventID}`)
      }
      throw new Error(`Failed to get scanner QR: ${response.status} ${response.statusText}`)
    }

    return await response.json()
  } catch (error) {
    if (error.message.includes('Failed to get scanner QR') || error.message.includes('Event not found')) {
      throw error
    }
    throw new Error(`Network error: ${error.message}`)
  }
}
