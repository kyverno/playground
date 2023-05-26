<template>
  <v-alert color="error" variant="outlined" class="mb-2" v-if="error">Failed to load namespaces: {{ error }}</v-alert>
  <v-autocomplete
    clearable
    variant="outlined"
    density="comfortable"
    hide-details
    :loading="loading"
    :items="data"
    @update:modelValue="(ns: string) => emit('update:modelValue', ns)"
    :model-value="prop.modelValue"
    label="Namespace" />
</template>

<script setup lang="ts">
import { useAPI } from '@/composables/api'

const prop = defineProps({
  modelValue: { type: String }
})

const emit = defineEmits(['update:modelValue'])

const { namespaces, error, loading, data } = useAPI<string[]>()

loading.value = true

namespaces()
  .then((ns) => {
    data.value = ns
  })
  .catch((err) => (error.value = err))
  .finally(() => (loading.value = false))
</script>
