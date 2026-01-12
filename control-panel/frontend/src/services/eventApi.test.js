import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createEvent, listEvents, getScannerQR } from './eventApi'

global.fetch = vi.fn()

describe('eventApi', () => {
  beforeEach(() => {
    fetch.mockClear()
  })

  describe('createEvent', () => {
    it('should call issuer service create-event endpoint with correct method', async () => {
      const mockResponse = {
        eventID: 'test-event',
        publicKey: 'test-public-key',
        qrCode: 'test-qr-code'
      }

      fetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse
      })

      await createEvent('test-event')

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/create-event',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Content-Type': 'application/json'
          })
        })
      )
    })
  })

  describe('listEvents', () => {
    it('should call issuer service list-events endpoint with GET method', async () => {
      const mockResponse = {
        events: [
          { eventID: 'event-1', publicKey: 'key-1' },
          { eventID: 'event-2', publicKey: 'key-2' }
        ]
      }

      fetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse
      })

      await listEvents()

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/list-events',
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            'Content-Type': 'application/json'
          })
        })
      )
    })

    it('should return events array from response', async () => {
      const mockResponse = {
        events: [
          { eventID: 'event-1', publicKey: 'key-1' },
          { eventID: 'event-2', publicKey: 'key-2' }
        ]
      }

      fetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse
      })

      const result = await listEvents()

      expect(result.events).toHaveLength(2)
      expect(result.events[0].eventID).toBe('event-1')
      expect(result.events[1].eventID).toBe('event-2')
    })

    it('should throw error when API request fails', async () => {
      fetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error'
      })

      await expect(listEvents()).rejects.toThrow('Failed to list events: 500 Internal Server Error')
    })

    it('should throw error when network request fails', async () => {
      fetch.mockRejectedValueOnce(new Error('Network error'))

      await expect(listEvents()).rejects.toThrow('Network error')
    })
  })

  describe('getScannerQR', () => {
    it('should call issuer service get-scanner-qr endpoint with GET method', async () => {
      const mockResponse = {
        eventID: 'test-event',
        publicKey: 'test-public-key',
        qrCode: 'test-qr-code-base64'
      }

      fetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse
      })

      await getScannerQR('test-event')

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/get-scanner-qr/test-event',
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            'Content-Type': 'application/json'
          })
        })
      )
    })

    it('should return scanner QR response with eventID, publicKey, and qrCode', async () => {
      const mockResponse = {
        eventID: 'test-event',
        publicKey: 'test-public-key',
        qrCode: 'test-qr-code-base64'
      }

      fetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse
      })

      const result = await getScannerQR('test-event')

      expect(result).toEqual(mockResponse)
      expect(result.eventID).toBe('test-event')
      expect(result.publicKey).toBe('test-public-key')
      expect(result.qrCode).toBe('test-qr-code-base64')
    })

    it('should throw error when event not found', async () => {
      fetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        statusText: 'Not Found'
      })

      await expect(getScannerQR('non-existent')).rejects.toThrow('Event not found: non-existent')
    })

    it('should throw error when API request fails', async () => {
      fetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error'
      })

      await expect(getScannerQR('test-event')).rejects.toThrow('Failed to get scanner QR: 500 Internal Server Error')
    })

    it('should throw error when network request fails', async () => {
      fetch.mockRejectedValueOnce(new Error('Network error'))

      await expect(getScannerQR('test-event')).rejects.toThrow('Network error')
    })
  })
})
