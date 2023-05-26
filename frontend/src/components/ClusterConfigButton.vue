<template>
  <v-dialog v-model="dialog" width="600px" :theme="layoutTheme">
    <template v-slot:activator="{ props }">
      <v-btn v-bind="props" prepend-icon="mdi-kubernetes" variant="flat" class="rounded-0 mx-0" color="success" width="25%">From Cluster</v-btn>
    </template>

    <v-card :theme="layoutTheme" title="Load from connected Cluster">
      <v-divider class="my-2" />
      <v-container>
        <simple-row v-if="error">
          <v-alert color="error" variant="outlined">{{ error }}</v-alert>
        </simple-row>
        <simple-row>
          <namespace-select v-model="namespace" />
        </simple-row>
        <simple-row>
          <v-text-field v-model="name" label="ConfigMap Name" variant="outlined" hide-details density="comfortable" />
        </simple-row>
      </v-container>
      <v-card-actions>
        <v-btn @click="dialog = false">Close</v-btn>
        <v-spacer />
        <v-btn :loading="loading" :disabled="!name || !namespace" @click="() => load()" :color="error ? 'error' : undefined">Load Config</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script setup lang="ts">
import { ref, watch } from 'vue'
import { layoutTheme } from '@/config'
import { useAPI, resourceToYAML } from '@/composables/api'
import { inputs } from '@/store'
import NamespaceSelect from '@/components/NamespaceSelect.vue'
import SimpleRow from '@/components/SimpleRow.vue'

const dialog = ref<boolean>(false)

const namespace = ref<string>('kyverno')
const name = ref<string>('kyverno')

const { loading, error, resource: loadResource } = useAPI<object>()

const load = () => {
  loading.value = true
  loadResource({ apiVersion: 'v1', kind: 'ConfigMap', namespace: namespace.value, name: name.value })
    .then((response) => {
      const results = resourceToYAML(response)

      inputs.config = results

      dialog.value = false
    })
    .catch((err) => {
      error.value = err
    })
    .finally(() => {
      loading.value = false
    })
}

watch(dialog, (value: boolean) => {
  if (value) return

  setTimeout(() => {
    loading.value = false
    error.value = undefined
  }, 300)
})
</script>
