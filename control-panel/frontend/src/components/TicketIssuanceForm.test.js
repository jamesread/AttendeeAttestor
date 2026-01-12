import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import TicketIssuanceForm from './TicketIssuanceForm.vue'
import { generateTicket } from '../services/ticketApi'

vi.mock('../services/ticketApi')

describe('TicketIssuanceForm', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should render form with eventID input field', () => {
    const wrapper = mount(TicketIssuanceForm)
    const eventIdInput = wrapper.find('input[name="eventID"]')
    expect(eventIdInput.exists()).toBe(true)
  })

  it('should render form with ticketID input field', () => {
    const wrapper = mount(TicketIssuanceForm)
    const ticketIdInput = wrapper.find('input[name="ticketID"]')
    expect(ticketIdInput.exists()).toBe(true)
  })

  it('should render form with ticketType input field', () => {
    const wrapper = mount(TicketIssuanceForm)
    const ticketTypeInput = wrapper.find('input[name="ticketType"]')
    expect(ticketTypeInput.exists()).toBe(true)
  })

  it('should render submit button', () => {
    const wrapper = mount(TicketIssuanceForm)
    const submitButton = wrapper.find('button[type="submit"]')
    expect(submitButton.exists()).toBe(true)
  })

  it('should call generateTicket when form is submitted', async () => {
    const mockResponse = {
      qrCode: 'test-qr-code',
      data: { eventID: 'test-event' },
      signature: 'test-signature'
    }

    generateTicket.mockResolvedValueOnce(mockResponse)

    const wrapper = mount(TicketIssuanceForm)

    await wrapper.find('input[name="eventID"]').setValue('test-event')
    await wrapper.find('input[name="ticketID"]').setValue('ticket-123')
    await wrapper.find('form').trigger('submit.prevent')

    await wrapper.vm.$nextTick()

    expect(generateTicket).toHaveBeenCalledWith({
      eventID: 'test-event',
      ticketID: 'ticket-123',
      ticketType: ''
    })
  })

  it('should display QR code when ticket is generated successfully', async () => {
    const mockResponse = {
      qrCode: 'test-qr-code-output',
      data: { eventID: 'test-event' },
      signature: 'test-signature'
    }

    generateTicket.mockResolvedValueOnce(mockResponse)

    const wrapper = mount(TicketIssuanceForm)

    await wrapper.find('input[name="eventID"]').setValue('test-event')
    await wrapper.find('form').trigger('submit.prevent')

    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(wrapper.text()).toContain('test-qr-code-output')
  })

  it('should display error message when ticket generation fails', async () => {
    generateTicket.mockRejectedValueOnce(new Error('API error'))

    const wrapper = mount(TicketIssuanceForm)

    await wrapper.find('input[name="eventID"]').setValue('test-event')
    await wrapper.find('form').trigger('submit.prevent')

    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(wrapper.text()).toContain('Failed to generate ticket')
  })

  it('should clear form after successful submission', async () => {
    const mockResponse = {
      qrCode: 'test-qr-code',
      data: { eventID: 'test-event' },
      signature: 'test-signature'
    }

    generateTicket.mockResolvedValueOnce(mockResponse)

    const wrapper = mount(TicketIssuanceForm)

    await wrapper.find('input[name="eventID"]').setValue('test-event')
    await wrapper.find('input[name="ticketID"]').setValue('ticket-123')
    await wrapper.find('form').trigger('submit.prevent')

    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    const eventIdInput = wrapper.find('input[name="eventID"]')
    expect(eventIdInput.element.value).toBe('')
  })
})

