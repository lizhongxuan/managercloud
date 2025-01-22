import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import MiddlewareList from '../views/middleware/List.vue'
import Monitor from '../views/Monitor.vue'

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard
  },
  {
    path: '/middleware',
    name: 'Middleware',
    component: MiddlewareList
  },
  {
    path: '/monitor',
    name: 'Monitor',
    component: Monitor
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router 