<template>
  <v-btn id="start-btn" size="large" :prepend-icon="icon" :color="color" :loading="loading" class="play" rounded @click="submit">Start</v-btn>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { MarkerSeverity, editor } from 'monaco-editor'
import { EngineResponse, ErrorResponse } from '@/types'
import { resolveAPI } from '@/utils'
import { inputs } from '@/store'

const props = defineProps({
  errorState: { type: Boolean, default: false }
})

const emit = defineEmits(['on-response', 'on-error'])

const loading = ref<boolean>(false)
const color = ref<string | undefined>('primary')
const icon = ref<string | undefined>('mdi-play')

watch(props, ({ errorState }: { errorState: boolean }) => {
  if (errorState) {
    color.value = 'error'
    icon.value = 'mdi-alert-circle-outline'
    return
  }

  color.value = 'primary'
  icon.value = 'mdi-play'
})

const api: string = resolveAPI()

const handleEditorErrors = () => {
  const markers = editor
    .getModelMarkers({ owner: 'yaml' })
    .filter((m) => m.severity == MarkerSeverity.Error)
    .map((m) => `${m.resource.path} L${m.startColumn}:${m.startLineNumber}: ${m.message}`)

  return markers
}

const submit = (): void => {
  const errors = handleEditorErrors()
  if (errors.length) {
    // emit("on-error", new Error(`<b>YAML validation failed, check the errors below.</b><br />${errors.join('<br />')}`));
    emit('on-error', new Error(`YAML validation failed, please check your manifests.`))
    return
  }

  if (!(inputs.policy || '').trim()) {
    emit('on-error', new Error('Policy is required'))
    return
  }

  if (!(inputs.resource || '').trim()) {
    emit('on-error', new Error('Resource is required'))
    return
  }

  loading.value = true

  fetch(`${api}/engine`, {
    body: JSON.stringify({
      policies: inputs.policy,
      resources: inputs.resource,
      oldResources: inputs.oldResource,
      policyExceptions: inputs.exceptions,
      clusterResources: inputs.clusterResources,
      context: inputs.context,
      config: inputs.config,
      customResourceDefinitions: inputs.customResourceDefinitions,
      imageData: inputs.imageData
    }),
    method: 'POST',
    mode: 'cors',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'application/json'
    }
  })
    .then((resp) => {
      if (resp.status === 400) {
        resp.json().then((err: ErrorResponse) => {
          if (err.violations) {
            emit(
              'on-error',
              new Error(
                err.violations.reduce((error, e) => {
                  const policy = [e.policyNamespace, e.policyName].filter((s) => !!s).join('/')

                  return error + `<h2 class="text-subtitle-2 mb-1">${policy}</h2><p class="mb-0 pb-0">${e.field}: ${e.detail}</p>`
                }, `<h1 class="text-subtitle-1 mb-2">Invalid Policy:</h1>`)
              )
            )
            return
          }

          emit('on-error', new Error(`ServerError: ${err.error}`))
        })
        return
      }

      if (resp.status > 300) {
        resp.text().then((err) => emit('on-error', new Error(`ServerError: ${err}`)))
        return
      }

      return resp
        .json()
        .catch(() => ({}))
        .then((content: EngineResponse) => {
          emit('on-response', content)
        })
    })
    .catch((err) => emit('on-error', err))
    .finally(() => (loading.value = false))
}
</script>

<style scoped>
.play {
  position: fixed;
  bottom: 45px;
  right: 50px;
}
</style>
