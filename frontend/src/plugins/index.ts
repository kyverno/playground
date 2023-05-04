import { loadFonts } from './webfontloader'
import vuetify from './vuetify'

import type { App } from 'vue'

export function registerPlugins (app: App) {
  loadFonts()
  app
    .use(vuetify)
}
