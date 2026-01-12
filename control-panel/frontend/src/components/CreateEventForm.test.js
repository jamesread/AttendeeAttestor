import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import CreateEventForm from './CreateEventForm.vue'
import { createEvent } from '../services/eventApi'

vi.mock('../services/eventApi')

describe('CreateEventForm', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should render form with eventID input field', () => {
    const wrapper = mount(CreateEventForm)
    const eventIdInput = wrapper.find('input[name="eventID"]')
    expect(eventIdInput.exists()).toBe(true)
  })

  it('should render submit button', () => {
    const wrapper = mount(CreateEventForm)
    const submitButton = wrapper.find('button[type="submit"]')
    expect(submitButton.exists()).toBe(true)
  })

  it('should call createEvent when form is submitted', async () => {
    const mockResponse = {
      eventID: 'test-event',
      publicKey: 'test-public-key',
      qrCode: 'test-qr-code'
    }

    createEvent.mockResolvedValueOnce(mockResponse)

    const wrapper = mount(CreateEventForm)

    await wrapper.find('input[name="eventID"]').setValue('test-event')
    await wrapper.find('form').trigger('submit.prevent')

    await wrapper.vm.$nextTick()

    expect(createEvent).toHaveBeenCalledWith('test-event')
  })

  it('should display QR code when event is created successfully', async () => {
    const mockResponse = {
      eventID: 'test-event',
      publicKey: 'test-public-key',
      qrCode: 'test-qr-code-output'
    }

    createEvent.mockResolvedValueOnce(mockResponse)

    const wrapper = mount(CreateEventForm)

    await wrapper.find('input[name="eventID"]').setValue('test-event')
    await wrapper.find('form').trigger('submit.prevent')

    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(wrapper.text()).toContain('test-qr-code-output')
  })

  it('should display public key when event is created successfully', async () => {
    const mockResponse = {
      eventID: 'test-event',
      publicKey: 'test-public-key-base64',
      qrCode: 'test-qr-code'
    }

    createEvent.mockResolvedValueOnce(mockResponse)

    const wrapper = mount(CreateEventForm)

    await wrapper.find('input[name="eventID"]').setValue('test-event')
    await wrapper.find('form').trigger('submit.prevent')

    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(wrapper.text()).toContain('test-public-key-base64')
  })

  it('should display error message when event creation fails', async () => {
    createEvent.mockRejectedValueOnce(new Error('API error'))

    const wrapper = mount(CreateEventForm)

    await wrapper.find('input[name="eventID"]').setValue('test-event')
    await wrapper.find('form').trigger('submit.prevent')

    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(wrapper.text()).toContain('Failed to create event')
  })

  it('should clear form after successful submission', async () => {
    const mockResponse = {
      eventID: 'test-event',
      publicKey: 'test-public-key',
      qrCode: 'test-qr-code'
    }

    createEvent.mockResolvedValueOnce(mockResponse)

    const wrapper = mount(CreateEventForm)

    await wrapper.find('input[name="eventID"]').setValue('test-event')
    await wrapper.find('form').trigger('submit.prevent')

    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    const eventIdInput = wrapper.find('input[name="eventID"]')
    expect(eventIdInput.element.value).toBe('')
  })
})

