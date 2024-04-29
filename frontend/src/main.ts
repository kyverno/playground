import App from './App.vue'
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import yamlWorker from './worker/yaml.worker.js?worker?worker'
import { setDiagnosticsOptions } from 'monaco-yaml'

import { createApp } from 'vue'
import { registerPlugins } from '@/plugins'

import policyexception from './schemas/policyexception-kyverno.io-v2beta1.json'
import policyexceptionv2 from './schemas/policyexception-kyverno.io-v2.json'
import clusterpolicyv1 from './schemas/clusterpolicy-kyverno.io-v1.json'
import policyv1 from './schemas/policy-kyverno.io-v1.json'
import clusterpolicyv2beta1 from './schemas/clusterpolicy-kyverno.io-v2beta1.json'
import policyv2beta1 from './schemas/policy-kyverno.io-v2beta1.json'
import vapv1beta1 from './schemas/validatingadmissionpolicy-admissionregistration-v1beta1.json'
import context from './schemas/context.json'
import { JSONSchema6 } from 'json-schema'

const baseURL = `${window.location.protocol}//${window.location.host}`

setDiagnosticsOptions({
  enableSchemaRequest: true,
  hover: true,
  completion: true,
  validate: true,
  format: true,
  schemas: [
    {
      schema: {
        oneOf: [
          { $ref: '#/definitions/clusterpolicy-v1' },
          { $ref: '#/definitions/policy-v1' },
          { $ref: '#/definitions/clusterpolicy-v2beta1' },
          { $ref: '#/definitions/policy-v2beta1' },
          { $ref: '#/definitions/vap-v1beta1' }
        ],
        definitions: {
          'clusterpolicy-v1': clusterpolicyv1 as JSONSchema6,
          'policy-v1': policyv1 as JSONSchema6,
          'clusterpolicy-v2beta1': clusterpolicyv2beta1 as JSONSchema6,
          'policy-v2beta1': policyv2beta1 as JSONSchema6,
          'vap-v1beta1': vapv1beta1 as JSONSchema6
        }
      },
      uri: `${baseURL}/schemas/policies.json`,
      fileMatch: ['policy.yaml']
    },
    { schema: policyexception as JSONSchema6, uri: `${baseURL}/schemas/policyexception-kyverno.io-v2beta1.json`, fileMatch: ['policyexception.yaml'] },
    { schema: policyexceptionv2 as JSONSchema6, uri: `${baseURL}/schemas/policyexception-kyverno.io-v2.json`, fileMatch: ['policyexception.yaml'] },
    { schema: context as JSONSchema6, uri: `${baseURL}/schemas/context.json`, fileMatch: ['context.yaml'] }
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
