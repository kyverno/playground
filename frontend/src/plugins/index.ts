import { loadFonts } from './webfontloader'
import vuetify from './vuetify'

import type { App } from 'vue'
import router from '@/router'

export function registerPlugins(app: App) {
  loadFonts()
  app.use(vuetify).use(router)
}
