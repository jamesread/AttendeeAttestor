import { createRouter, createWebHistory } from 'vue-router'
import CreateEventForm from './components/CreateEventForm.vue'
import TicketIssuanceForm from './components/TicketIssuanceForm.vue'
import EventList from './components/EventList.vue'
import IssuedTicketsList from './components/IssuedTicketsList.vue'
import CBORDiagnostic from './components/CBORDiagnostic.vue'
import TicketDetails from './components/TicketDetails.vue'
import Scanner from './components/Scanner.vue'
import Admin from './components/Admin.vue'

const routes = [
  {
    path: '/',
    redirect: '/scanner'
  },
  {
    path: '/admin',
    name: 'admin',
    component: Admin
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
  },
  {
    path: '/issued-tickets',
    name: 'issued-tickets',
    component: IssuedTicketsList
  },
  {
    path: '/diagnostic',
    name: 'diagnostic',
    component: CBORDiagnostic
  },
  {
    path: '/ticket-details',
    name: 'ticket-details',
    component: TicketDetails
  },
  {
    path: '/scanner',
    name: 'scanner',
    component: Scanner
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router

