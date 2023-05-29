import App from './App.vue'
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import yamlWorker from './worker/yaml.worker.js?worker?worker'
import { setDiagnosticsOptions } from 'monaco-yaml'

import { createApp } from 'vue'
import { registerPlugins } from '@/plugins'

import clusterpolicy from './schemas/clusterpolicy.json'
import context from './schemas/context.json'

const baseURL = `${window.location.protocol}//${window.location.host}`

setDiagnosticsOptions({
  enableSchemaRequest: true,
  hover: true,
  completion: true,
  validate: true,
  format: true,
  schemas: [
    { ...clusterpolicy, uri: `${baseURL}/schemas/clusterpolicy.json`, fileMatch: ['policy.yaml'] },
    { ...context, uri: `${baseURL}/schemas/context.json`, fileMatch: ['context.yaml'] },
    { ...context, uri: `${baseURL}/schemas/policyexception.json`, fileMatch: ['policyexception.yaml'] }
  ]
})

// @ts-ignore
self.MonacoEnvironment = {
  getWorker(_: any, label: string) {
    if (label === 'yaml') {
      return new yamlWorker()
    }
    if (label === 'json') {
      return new jsonWorker()
    }
    return new editorWorker()
  }
}

const app = createApp(App)

registerPlugins(app)

app.mount('#app')
