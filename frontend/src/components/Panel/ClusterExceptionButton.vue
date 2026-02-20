<template>
  <cluster-button
    :default-resource="defaultResource"
    :model-value="modelValue"
    @update:model-value="(res: string) => emit('update:modelValue', res)"
    label="From Cluster"
  />
</template>
<script setup lang="ts">
import { ref } from 'vue'
import ClusterButton from '../ClusterButton.vue'

const defaultResource = ref({
  title: 'kyverno.io/v2alpha1 PolicyException',
  apiVersion: 'kyverno.io/v2alpha1',
  kind: 'PolicyException',
  clusterScoped: false,
})

defineProps({
  modelValue: { type: String },
})

const emit = defineEmits(['update:modelValue'])
</script>
