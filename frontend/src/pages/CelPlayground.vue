<template>
  <v-app :theme="layoutTheme">
    <AppBar>
      <template #append-actions>
        <v-btn to="/" prepend-icon="mdi-arrow-left" class="d-md-flex d-none">Back to Playground</v-btn>
      </template>
      <template #mobile-actions>
        <v-btn to="/" icon="mdi-arrow-left" />
      </template>
    </AppBar>
    <v-main>
      <v-container fluid class="pr-lg-4 pr-md-4 pr-sm-8 pr-8">
        <v-row>
          <v-col :md="6" :sm="12">
            <v-card>
              <v-toolbar color="#3783c4" theme="dark" density="compact" flat>
                <v-toolbar-title>Expression</v-toolbar-title>
                <v-toolbar-items>
                  <CopyButton :value="expression" />
                </v-toolbar-items>
              </v-toolbar>
              <MonacoEditor
                id="cel-expression"
                language="plaintext"
                :theme="editorTheme"
                v-model="expression"
                :options="expressionOptions"
                style="height: 120px"
              />
            </v-card>

            <v-card class="mt-4">
              <EditorToolbar title="Resource" info="Bound to the `object` CEL variable" v-model="resource">
                <template #prepend-actions>
                  <TemplateButton @select="(template: string) => (resource = template)" />
                </template>
              </EditorToolbar>
              <MonacoEditor
                id="cel-resource"
                language="yaml"
                :theme="editorTheme"
                v-model="resource"
                :options="resourceOptions"
                :uri="resourceUri"
                style="height: 280px"
              />
            </v-card>

            <v-card class="mt-4">
              <EditorToolbar
                title="Old Resource"
                info="Bound to the `oldObject` CEL variable"
                v-model="oldResource"
              >
                <template #append-actions>
                  <v-btn :icon="collapseOld ? 'mdi-chevron-down' : 'mdi-chevron-up'" @click="collapseOld = !collapseOld" />
                </template>
              </EditorToolbar>
              <MonacoEditor
                id="cel-old-resource"
                language="yaml"
                :theme="editorTheme"
                v-model="oldResource"
                :options="resourceOptions"
                :uri="oldResourceUri"
                style="height: 200px"
                v-show="!collapseOld"
              />
            </v-card>

            <v-card class="mt-4">
              <v-toolbar color="#3783c4" theme="dark" density="compact" flat>
                <v-toolbar-title>
                  Context
                  <v-tooltip
                    text="Used to build the `request` CEL variable"
                    content-class="no-opacity-tooltip"
                  >
                    <template v-slot:activator="{ props }">
                      <v-btn v-bind="props" icon="mdi-information-outline" variant="text" size="small" />
                    </template>
                  </v-tooltip>
                </v-toolbar-title>
                <v-toolbar-items>
                  <v-btn :icon="collapseContext ? 'mdi-chevron-down' : 'mdi-chevron-up'" @click="collapseContext = !collapseContext" />
                </v-toolbar-items>
              </v-toolbar>
              <v-card-text v-show="!collapseContext">
                <v-row>
                  <v-col cols="12" sm="6">
                    <v-select
                      label="Operation"
                      :items="operations"
                      v-model="operation"
                      variant="outlined"
                      density="comfortable"
                      hide-details
                    />
                  </v-col>
                  <v-col cols="12" sm="6">
                    <v-text-field
                      label="Username"
                      v-model="username"
                      variant="outlined"
                      density="comfortable"
                      hide-details
                    />
                  </v-col>
                  <v-col cols="12">
                    <v-text-field
                      label="Groups (comma separated)"
                      v-model="groupsInput"
                      variant="outlined"
                      density="comfortable"
                      hide-details
                    />
                  </v-col>
                </v-row>
              </v-card-text>
            </v-card>
          </v-col>

          <v-col :md="6" :sm="12">
            <v-card>
              <v-toolbar color="#3783c4" theme="dark" density="compact" flat>
                <v-toolbar-title>Result</v-toolbar-title>
                <v-toolbar-items v-if="result?.type">
                  <v-chip class="align-self-center mr-4" size="small" color="white" variant="outlined">
                    {{ result.type }}
                  </v-chip>
                </v-toolbar-items>
              </v-toolbar>
              <v-card-text>
                <v-alert v-if="error" type="error" variant="tonal" class="mb-0">{{ error }}</v-alert>
                <v-alert v-else-if="result?.error" type="warning" variant="tonal" class="mb-0">
                  {{ result.error }}
                </v-alert>
                <pre v-else-if="result" class="result-value">{{ formattedValue }}</pre>
                <div v-else class="text-medium-emphasis pa-4">
                  Evaluate an expression to see its result here.
                </div>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <v-btn
          id="evaluate-btn"
          size="large"
          prepend-icon="mdi-play"
          color="primary"
          class="evaluate"
          rounded
          :loading="loading"
          @click="submit"
        >
          Evaluate
        </v-btn>
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Uri } from 'monaco-editor'
import { editorTheme, layoutTheme } from '@/config'
import { AppBar } from '@/components/AppBar'
import MonacoEditor from '@/components/Panel/MonacoEditor.vue'
import EditorToolbar from '@/components/Panel/EditorToolbar.vue'
import CopyButton from '@/components/Panel/CopyButton.vue'
import TemplateButton from '@/components/Panel/TemplateButton.vue'
import { useCelAPI } from '@/composables/cel'
import type { CelEvaluateResponse } from '@/types'

const expression = ref<string>('')
const resource = ref<string>('')
const oldResource = ref<string>('')
const collapseOld = ref<boolean>(true)
const collapseContext = ref<boolean>(true)

const operations = ['CREATE', 'UPDATE', 'DELETE', 'CONNECT']
const operation = ref<string>('CREATE')
const username = ref<string>('')
const groupsInput = ref<string>('')

const expressionOptions = {
  colorDecorators: true,
  lineHeight: 24,
  tabSize: 2,
  minimap: { enabled: false },
  lineNumbers: 'off' as const,
}

// Distinct URIs from the Home page's `resource.yaml` model (see ResourceEditor.vue) -
// Monaco models are looked up globally by URI, so reusing that URI here would bind
// this editor to whatever content the Home page's resource editor last held.
const resourceUri = Uri.parse('cel-resource.yaml')
const oldResourceUri = Uri.parse('cel-old-resource.yaml')
const resourceOptions = {
  colorDecorators: true,
  lineHeight: 24,
  tabSize: 2,
}

const { evaluate } = useCelAPI()

const loading = ref<boolean>(false)
const result = ref<CelEvaluateResponse>()
const error = ref<string>('')

const formattedValue = computed(() => {
  if (result.value?.value === undefined) return ''
  return JSON.stringify(result.value.value, null, 2)
})

const submit = () => {
  if (!expression.value.trim()) {
    error.value = 'Expression is required'
    result.value = undefined
    return
  }

  error.value = ''
  loading.value = true

  const groups = groupsInput.value
    .split(',')
    .map((g) => g.trim())
    .filter((g) => !!g)

  evaluate({
    expression: expression.value,
    resource: resource.value,
    oldResource: oldResource.value,
    context: {
      operation: operation.value,
      username: username.value,
      groups,
    },
  })
    .then((resp) => {
      result.value = resp
    })
    .catch((err: Error) => {
      error.value = err.message
      result.value = undefined
    })
    .finally(() => {
      loading.value = false
    })
}
</script>

<style scoped>
.evaluate {
  position: fixed;
  bottom: 45px;
  right: 50px;
}

.result-value {
  white-space: pre-wrap;
  word-break: break-word;
  font-family: monospace;
  margin: 0;
}
</style>
