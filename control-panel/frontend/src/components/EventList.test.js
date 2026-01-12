import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import EventList from './EventList.vue'
import { listEvents, getScannerQR } from '../services/eventApi'

vi.mock('../services/eventApi')

describe('EventList', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should render scanner setup button for each event', async () => {
    const mockEvents = [
      { eventID: 'event-1', publicKey: 'key-1' },
      { eventID: 'event-2', publicKey: 'key-2' }
    ]

    listEvents.mockResolvedValueOnce({ events: mockEvents })

    const wrapper = mount(EventList)
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    const buttons = wrapper.findAll('.scanner-setup-button')
    expect(buttons.length).toBe(2)
  })

  it('should show QR modal when scanner setup button is clicked', async () => {
    const mockEvents = [
      { eventID: 'event-1', publicKey: 'key-1' }
    ]

    listEvents.mockResolvedValueOnce({ events: mockEvents })
    getScannerQR.mockResolvedValueOnce({
      eventID: 'event-1',
      publicKey: 'key-1',
      qrCode: 'test-base64-qr-code'
    })

    const wrapper = mount(EventList)
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    const button = wrapper.find('.scanner-setup-button')
    await button.trigger('click')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(wrapper.find('.qr-modal-overlay').exists()).toBe(true)
  })

  it('should call getScannerQR API when button is clicked', async () => {
    const mockEvents = [
      { eventID: 'test-event', publicKey: 'test-key' }
    ]

    listEvents.mockResolvedValueOnce({ events: mockEvents })
    getScannerQR.mockResolvedValueOnce({
      eventID: 'test-event',
      publicKey: 'test-key',
      qrCode: 'test-base64-qr-code'
    })

    const wrapper = mount(EventList)
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    const button = wrapper.find('.scanner-setup-button')
    await button.trigger('click')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(getScannerQR).toHaveBeenCalledWith('test-event')
  })

  it('should close QR modal when close button is clicked', async () => {
    const mockEvents = [
      { eventID: 'event-1', publicKey: 'key-1' }
    ]

    listEvents.mockResolvedValueOnce({ events: mockEvents })
    getScannerQR.mockResolvedValueOnce({
      eventID: 'event-1',
      publicKey: 'key-1',
      qrCode: 'test-base64-qr-code'
    })

    const wrapper = mount(EventList)
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    const setupButton = wrapper.find('.scanner-setup-button')
    await setupButton.trigger('click')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    const closeButton = wrapper.find('.close-button')
    await closeButton.trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.qr-modal-overlay').exists()).toBe(false)
  })

  it('should close QR modal when overlay is clicked', async () => {
    const mockEvents = [
      { eventID: 'event-1', publicKey: 'key-1' }
    ]

    listEvents.mockResolvedValueOnce({ events: mockEvents })
    getScannerQR.mockResolvedValueOnce({
      eventID: 'event-1',
      publicKey: 'key-1',
      qrCode: 'test-base64-qr-code'
    })

    const wrapper = mount(EventList)
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    const setupButton = wrapper.find('.scanner-setup-button')
    await setupButton.trigger('click')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    const overlay = wrapper.find('.qr-modal-overlay')
    await overlay.trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.qr-modal-overlay').exists()).toBe(false)
  })

  it('should display event ID in modal header', async () => {
    const mockEvents = [
      { eventID: 'my-event', publicKey: 'my-key' }
    ]

    listEvents.mockResolvedValueOnce({ events: mockEvents })
    getScannerQR.mockResolvedValueOnce({
      eventID: 'my-event',
      publicKey: 'my-key',
      qrCode: 'test-base64-qr-code'
    })

    const wrapper = mount(EventList)
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    const setupButton = wrapper.find('.scanner-setup-button')
    await setupButton.trigger('click')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(wrapper.text()).toContain('my-event')
  })

  it('should display QR code image when loaded from API', async () => {
    const mockEvents = [
      { eventID: 'test-event', publicKey: 'test-key' }
    ]

    listEvents.mockResolvedValueOnce({ events: mockEvents })
    getScannerQR.mockResolvedValueOnce({
      eventID: 'test-event',
      publicKey: 'test-key',
      qrCode: 'test-base64-data'
    })

    const wrapper = mount(EventList)
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    const setupButton = wrapper.find('.scanner-setup-button')
    await setupButton.trigger('click')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    const qrImage = wrapper.find('.qr-code-image')
    expect(qrImage.exists()).toBe(true)
    expect(qrImage.attributes('src')).toContain('test-base64-data')
  })

  it('should display error message when getScannerQR fails', async () => {
    const mockEvents = [
      { eventID: 'test-event', publicKey: 'test-key' }
    ]

    listEvents.mockResolvedValueOnce({ events: mockEvents })
    getScannerQR.mockRejectedValueOnce(new Error('API error'))

    const wrapper = mount(EventList)
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    const setupButton = wrapper.find('.scanner-setup-button')
    await setupButton.trigger('click')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(wrapper.text()).toContain('Failed to load QR code')
  })
})
