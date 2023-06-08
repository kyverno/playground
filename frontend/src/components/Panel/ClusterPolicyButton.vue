<template>
  <cluster-button :default-resource="defaultResource" v-model="inputs.policy">
    <template #resource-api="{ update, resource }">
      <simple-row>
        <policy-type-select :model-value="resource" @update:modelValue="update" />
      </simple-row>
    </template>
  </cluster-button>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import ClusterButton from '../ClusterButton.vue'
import PolicyTypeSelect from '../PolicyTypeSelect.vue'
import SimpleRow from '../SimpleRow.vue'
import { inputs } from '@/store'

const defaultResource = ref({ title: 'kyverno.io/v1 ClusterPolicy', apiVersion: 'kyverno.io/v1', kind: 'ClusterPolicy', clusterScoped: true })
</script>
