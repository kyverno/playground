<template>
  <cluster-button
    :default-resource="defaultResource"
    :model-value="modelValue"
    @update:model-value="(res: string) => emit('update:modelValue', res)"
  >
    <template #resource-api="{ update, resource }">
      <simple-row>
        <resource-type-select :model-value="resource" @update:modelValue="update" />
      </simple-row>
    </template>
  </cluster-button>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import ClusterButton from '../ClusterButton.vue'
import SimpleRow from '../SimpleRow.vue'
import ResourceTypeSelect from '../ResourceTypeSelect.vue'

const defaultResource = ref({
  apiVersion: 'v1',
  kind: 'Pod',
  title: 'v1 Pod',
  clusterScoped: false,
})

defineProps({
  modelValue: { type: String },
})

const emit = defineEmits(['update:modelValue'])
</script>
