import Home from '@/pages/Home.vue'
import MutationDetails from '@/pages/MutationDetails.vue'
import { RouteRecordRaw, createRouter, createWebHashHistory } from 'vue-router'

const routes: RouteRecordRaw[] = [{
  path: '/mutation/:id',
  component: MutationDetails,
  name: 'mutation-details'
}, {
  path: '/',
  component: Home,
  name: 'home'
}]

const router = createRouter({
  history: createWebHashHistory(process.env.BASE_URL),
  routes,
})

export default router
