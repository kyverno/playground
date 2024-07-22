import Home from '@/pages/Home.vue'
import MutationDetails from '@/pages/MutationDetails.vue'
import GenerationDetails from '@/pages/GenerationDetails.vue'
import { type RouteRecordRaw, createRouter, createWebHashHistory } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/mutation/:id',
    component: MutationDetails,
    name: 'mutation-details'
  },
  {
    path: '/generation/:id',
    component: GenerationDetails,
    name: 'generation-details'
  },
  {
    path: '/',
    component: Home,
    name: 'home'
  }
]

const router = createRouter({
  history: createWebHashHistory(process.env.BASE_URL),
  routes
})

export default router
