// Plugins
import vue from '@vitejs/plugin-vue'
import vuetify, { transformAssetUrls } from 'vite-plugin-vuetify'
import vueDevTools from 'vite-plugin-vue-devtools'

// Utilities
import { defineConfig } from 'vite'
import { fileURLToPath, URL } from 'node:url'

// https://vitejs.dev/config/
export default defineConfig({
  base: '',
  plugins: [
    vue({
      template: { transformAssetUrls }
    }),
    // https://github.com/vuetifyjs/vuetify-loader/tree/next/packages/vite-plugin
    vuetify({
      autoImport: true
    }),
    vueDevTools()
  ],
  define: {
    'process.env': {},
    APP_VERSION: JSON.stringify(process.env.APP_VERSION || '')
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      // monaco-editor >= 0.56 restricts deep imports via its exports map,
      // but monaco-worker-manager (used by monaco-yaml) still uses the old path
      'monaco-editor/esm/vs/editor/editor.worker.js': 'monaco-editor/editor/editor.worker.js'
    },
    extensions: ['.js', '.json', '.jsx', '.mjs', '.ts', '.tsx', '.vue']
  },
  server: {
    port: 3000
  }
})
