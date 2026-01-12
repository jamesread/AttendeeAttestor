const ISSUER_SERVICE_URL = 'http://localhost:8080/generate'
const LIST_ISSUED_TICKETS_URL = 'http://localhost:8080/issued-tickets'
const GET_TICKET_BY_ID_URL = 'http://localhost:8080/get-ticket-by-id'

export async function generateTicket(ticketData) {
  try {
    const response = await fetch(ISSUER_SERVICE_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(ticketData)
    })

    if (!response.ok) {
      throw new Error(`Failed to generate ticket: ${response.status} ${response.statusText}`)
    }

    return await response.json()
  } catch (error) {
    if (error.message.includes('Failed to generate ticket')) {
      throw error
    }
    throw new Error(`Network error: ${error.message}`)
  }
}

export async function listIssuedTickets(eventID) {
  try {
    const url = eventID 
      ? `${LIST_ISSUED_TICKETS_URL}?eventID=${encodeURIComponent(eventID)}`
      : LIST_ISSUED_TICKETS_URL
    
    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      throw new Error(`Failed to list issued tickets: ${response.status} ${response.statusText}`)
    }

    return await response.json()
  } catch (error) {
    if (error.message.includes('Failed to list issued tickets')) {
      throw error
    }
    throw new Error(`Network error: ${error.message}`)
  }
}

export async function getTicketByID(ticketID) {
  try {
    const url = `${GET_TICKET_BY_ID_URL}?ticketID=${encodeURIComponent(ticketID)}`
    
    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      throw new Error(errorData.error || `Failed to get ticket: ${response.status} ${response.statusText}`)
    }

    const result = await response.json()
    if (result.error) {
      throw new Error(result.error)
    }

    return result
  } catch (error) {
    if (error.message.includes('Failed to get ticket') || error.message.includes('Ticket not found')) {
      throw error
    }
    throw new Error(`Network error: ${error.message}`)
  }
}

