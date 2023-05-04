import { RouteRecordRaw, createRouter, createWebHistory } from 'vue-router'

const routes: RouteRecordRaw[] = []

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
})

export default router
