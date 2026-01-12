import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import CreateEventForm from './components/CreateEventForm.vue'
import TicketIssuanceForm from './components/TicketIssuanceForm.vue'
import EventList from './components/EventList.vue'

const routes = [
  {
    path: '/',
    redirect: '/list-events'
  },
  {
    path: '/create-event',
    name: 'create-event',
    component: CreateEventForm
  },
  {
    path: '/list-events',
    name: 'list-events',
    component: EventList
  },
  {
    path: '/issue-ticket',
    name: 'issue-ticket',
    component: TicketIssuanceForm
  }
]

function createTestRouter() {
  return createRouter({
    history: createWebHistory(),
    routes
  })
}

describe('App', () => {
  it('should render navigation with router links', async () => {
    const router = createTestRouter()
    const wrapper = mount(App, {
      global: {
        plugins: [router]
      }
    })
    await router.isReady()

    const navigation = wrapper.find('.navigation')
    expect(navigation.exists()).toBe(true)
  })

  it('should not render Create Event router link in navigation', async () => {
    const router = createTestRouter()
    const wrapper = mount(App, {
      global: {
        plugins: [router]
      }
    })
    await router.isReady()

    const createEventLink = wrapper.find('a[href="/create-event"]')
    expect(createEventLink.exists()).toBe(false)
  })

  it('should not render Issue Ticket router link in navigation', async () => {
    const router = createTestRouter()
    const wrapper = mount(App, {
      global: {
        plugins: [router]
      }
    })
    await router.isReady()

    const issueTicketLink = wrapper.find('a[href="/issue-ticket"]')
    expect(issueTicketLink.exists()).toBe(false)
  })

  it('should show CreateEventForm when on create-event route', async () => {
    const router = createTestRouter()
    await router.push('/create-event')
    await router.isReady()

    const wrapper = mount(App, {
      global: {
        plugins: [router]
      }
    })
    await router.isReady()

    const createEventForm = wrapper.findComponent({ name: 'CreateEventForm' })
    expect(createEventForm.exists()).toBe(true)
  })

  it('should show TicketIssuanceForm when on issue-ticket route', async () => {
    const router = createTestRouter()
    await router.push('/issue-ticket')
    await router.isReady()

    const wrapper = mount(App, {
      global: {
        plugins: [router]
      }
    })
    await router.isReady()

    const ticketForm = wrapper.findComponent({ name: 'TicketIssuanceForm' })
    expect(ticketForm.exists()).toBe(true)
  })

  it('should redirect root path to list-events', async () => {
    const router = createTestRouter()
    await router.push('/')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/list-events')
  })
})

