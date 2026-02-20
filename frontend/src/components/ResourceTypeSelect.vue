<template>
  <v-alert color="error" variant="outlined" class="mb-2" v-if="error"
    >Failed to load resource types: {{ error }}</v-alert
  >
  <v-autocomplete
    variant="outlined"
    hide-details
    :items="state.kinds"
    label="Resource Type"
    v-model="resource"
    return-object
    density="comfortable"
    :loading="loading"
  />
</template>

<script setup lang="ts">
import { type Resource, useAPI } from '@/composables/api'
import { computed } from 'vue'
import { type PropType } from 'vue'
import { state } from '@/store'

const { kinds, loading, error } = useAPI()

loading.value = true
kinds()
  .then((kinds) => {
    state.kinds = kinds.map((k) => ({ ...k, title: `${k.apiVersion} ${k.kind}` }))
  })
  .catch((err) => {
    error.value = err
  })
  .finally(() => {
    loading.value = false
  })

const props = defineProps({
  modelValue: { type: Object as PropType<Resource>, required: true },
})

const emit = defineEmits(['update:modelValue'])

const resource = computed({
  get() {
    return props.modelValue
  },
  set(val: Resource) {
    emit('update:modelValue', { ...val })
  },
})
</script>
