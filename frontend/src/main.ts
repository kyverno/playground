import App from './App.vue'
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import yamlWorker from './worker/yaml.worker.js?worker?worker'
import * as monaco from 'monaco-editor'
import { configureMonacoYaml, type JSONSchema } from 'monaco-yaml'

import { createApp } from 'vue'
import { registerPlugins } from '@/plugins'

import policyexception from './schemas/policyexception-kyverno.io-v2beta1.json'
import policyexceptionv2 from './schemas/policyexception-kyverno.io-v2.json'
import clusterpolicyv1 from './schemas/clusterpolicy-kyverno.io-v1.json'
import policyv1 from './schemas/policy-kyverno.io-v1.json'
import clusterpolicyv2beta1 from './schemas/clusterpolicy-kyverno.io-v2beta1.json'
import policyv2beta1 from './schemas/policy-kyverno.io-v2beta1.json'
import vapv1 from './schemas/validatingadmissionpolicy-admissionregistration-v1.json'
import vapbv1 from './schemas/validatingadmissionpolicybinding-admissionregistration-v1.json'
import vpolv1alpha1 from './schemas/validatingpolicy-policies.kyverno.io-v1alpha1.json'
import ivpolv1alpha1 from './schemas/imagevalidatingpolicy-policies.kyverno.io-v1alpha1.json'
import dpolv1alpha1 from './schemas/deletingpolicy-policies.kyverno.io-v1alpha1.json'
import gpolv1alpha1 from './schemas/generatingpolicy-policies.kyverno.io-v1alpha1.json'
import mpolv1alpha1 from './schemas/mutatingpolicy-policies.kyverno.io-v1alpha1.json'
import context from './schemas/context.json'

const baseURL = `${window.location.protocol}//${window.location.host}`

configureMonacoYaml(monaco, {
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
          { $ref: '#/definitions/vap-v1' },
          { $ref: '#/definitions/vapb-v1' },
          { $ref: '#/definitions/validatingpolicy-v1alpha1' },
          { $ref: '#/definitions/imagevalidatingpolicy-v1alpha1' },
          { $ref: '#/definitions/deletingpolicy-v1alpha1' },
          { $ref: '#/definitions/generatingpolicy-v1alpha1' },
          { $ref: '#/definitions/mutatingpolicy-v1alpha1' }
        ],
        definitions: {
          'clusterpolicy-v1': clusterpolicyv1 as JSONSchema,
          'policy-v1': policyv1 as JSONSchema,
          'clusterpolicy-v2beta1': clusterpolicyv2beta1 as JSONSchema,
          'policy-v2beta1': policyv2beta1 as JSONSchema,
          'vap-v1': vapv1 as JSONSchema,
          'vapb-v1': vapbv1 as JSONSchema,
          'validatingpolicy-v1alpha1': vpolv1alpha1 as JSONSchema,
          'imagevalidatingpolicy-v1alpha1': ivpolv1alpha1 as JSONSchema,
          'deletingpolicy-v1alpha1': dpolv1alpha1 as JSONSchema,
          'generatingpolicy-v1alpha1': gpolv1alpha1 as JSONSchema,
          'mutatingpolicy-v1alpha1': mpolv1alpha1 as JSONSchema
        }
      },
      uri: `${baseURL}/schemas/policies.json`,
      fileMatch: ['policy.yaml']
    },
    {
      schema: policyexception as JSONSchema,
      uri: `${baseURL}/schemas/policyexception-kyverno.io-v2beta1.json`,
      fileMatch: ['policyexception.yaml']
    },
    {
      schema: policyexceptionv2 as JSONSchema,
      uri: `${baseURL}/schemas/policyexception-kyverno.io-v2.json`,
      fileMatch: ['policyexception.yaml']
    },
    {
      schema: vpolv1alpha1 as JSONSchema,
      uri: `${baseURL}/schemas/validatingpolicy-policies.kyverno.io-v1alpha1.json`,
      fileMatch: ['validatingpolicy.yaml']
    },
    {
      schema: ivpolv1alpha1 as JSONSchema,
      uri: `${baseURL}/schemas/imagevalidatingpolicy-policies.kyverno.io-v1alpha1.json`,
      fileMatch: ['validatingpolicy.yaml']
    },
    {
      schema: dpolv1alpha1 as JSONSchema,
      uri: `${baseURL}/schemas/deletingpolicy-policies.kyverno.io-v1alpha1.json`,
      fileMatch: ['deletingpolicy.yaml']
    },
    {
      schema: mpolv1alpha1 as JSONSchema,
      uri: `${baseURL}/schemas/mutatingpolicy-policies.kyverno.io-v1alpha1.json`,
      fileMatch: ['mutatingpolicy.yaml']
    },
    {
      schema: gpolv1alpha1 as JSONSchema,
      uri: `${baseURL}/schemas/generatingpolicy-policies.kyverno.io-v1alpha1.json`,
      fileMatch: ['generatingpolicy.yaml']
    },
    {
      schema: context as JSONSchema,
      uri: `${baseURL}/schemas/context.json`,
      fileMatch: ['context.yaml']
    }
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
