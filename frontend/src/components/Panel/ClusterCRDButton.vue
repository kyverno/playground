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
  title: 'apiextensions.k8s.io/v1 CustomResourceDefinition',
  apiVersion: 'apiextensions.k8s.io/v1',
  kind: 'CustomResourceDefinition',
  clusterScoped: true
})

defineProps({
  modelValue: { type: String }
})

const emit = defineEmits(['update:modelValue'])
</script>
