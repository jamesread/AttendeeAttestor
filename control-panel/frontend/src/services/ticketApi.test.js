import { describe, it, expect, vi, beforeEach } from 'vitest'
import { generateTicket } from './ticketApi'

global.fetch = vi.fn()

describe('ticketApi', () => {
  beforeEach(() => {
    fetch.mockClear()
  })

  it('should call issuer service with correct endpoint and method', async () => {
    const mockResponse = {
      qrCode: 'test-qr-code',
      data: { eventID: 'test-event' },
      signature: 'test-signature'
    }

    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockResponse
    })

    const ticketData = { eventID: 'test-event', ticketID: 'ticket-123' }
    await generateTicket(ticketData)

    expect(fetch).toHaveBeenCalledWith(
      'http://localhost:8080/generate',
      expect.objectContaining({
        method: 'POST',
        headers: expect.objectContaining({
          'Content-Type': 'application/json'
        })
      })
    )
  })

  it('should send ticket data as JSON in request body', async () => {
    const mockResponse = {
      qrCode: 'test-qr-code',
      data: { eventID: 'test-event' },
      signature: 'test-signature'
    }

    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockResponse
    })

    const ticketData = {
      eventID: 'test-event',
      ticketID: 'ticket-123',
      ticketType: 'VIP'
    }

    await generateTicket(ticketData)

    const callArgs = fetch.mock.calls[0]
    const requestBody = JSON.parse(callArgs[1].body)

    expect(requestBody).toEqual(ticketData)
  })

  it('should return ticket response with qrCode, data, and signature', async () => {
    const mockResponse = {
      qrCode: 'test-qr-code',
      data: { eventID: 'test-event', ticketID: 'ticket-123' },
      signature: 'test-signature'
    }

    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockResponse
    })

    const ticketData = { eventID: 'test-event' }
    const result = await generateTicket(ticketData)

    expect(result).toEqual(mockResponse)
    expect(result.qrCode).toBe('test-qr-code')
    expect(result.data).toEqual({ eventID: 'test-event', ticketID: 'ticket-123' })
    expect(result.signature).toBe('test-signature')
  })

  it('should throw error when API request fails', async () => {
    fetch.mockResolvedValueOnce({
      ok: false,
      status: 400,
      statusText: 'Bad Request'
    })

    const ticketData = { eventID: 'test-event' }

    await expect(generateTicket(ticketData)).rejects.toThrow('Failed to generate ticket: 400 Bad Request')
  })

  it('should throw error when network request fails', async () => {
    fetch.mockRejectedValueOnce(new Error('Network error'))

    const ticketData = { eventID: 'test-event' }

    await expect(generateTicket(ticketData)).rejects.toThrow('Network error')
  })
})

