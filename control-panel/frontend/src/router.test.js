import { describe, it, expect } from 'vitest'
import { createRouter, createWebHistory } from 'vue-router'
import { mount } from '@vue/test-utils'
import App from './App.vue'
import CreateEventForm from './components/CreateEventForm.vue'
import TicketIssuanceForm from './components/TicketIssuanceForm.vue'

const routes = [
  { path: '/', redirect: '/list-events' },
  { path: '/create-event', component: CreateEventForm },
  { path: '/list-events', component: () => import('./components/EventList.vue') },
  { path: '/issue-ticket', component: TicketIssuanceForm }
]

describe('Router', () => {
  it('should redirect root path to list-events', async () => {
    const router = createRouter({
      history: createWebHistory(),
      routes
    })

    await router.push('/')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/list-events')
  })

  it('should navigate to create-event route', async () => {
    const router = createRouter({
      history: createWebHistory(),
      routes
    })

    await router.push('/create-event')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/create-event')
  })

  it('should navigate to issue-ticket route', async () => {
    const router = createRouter({
      history: createWebHistory(),
      routes
    })

    await router.push('/issue-ticket')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/issue-ticket')
  })
})

