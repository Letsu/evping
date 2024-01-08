// Composables
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    component: () => import('@/layouts/default/Default.vue'),
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/Home.vue'),
      },
      {
        path: '',
        name: 'AllHosts',
        component: () => import('@/views/AllHosts.vue'),
      },
      {
        path: '/host/:host',
        name: 'Host',
        component: () => import('@/views/Host.vue'),
        props: true
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
})

export default router
