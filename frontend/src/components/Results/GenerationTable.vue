<template>
  <v-card-title class="my-2 text-subtitle-1">{{ title }}</v-card-title>
  <v-divider />
  <v-data-table
    density="default"
    hover
    :items="items"
    :headers="headers as any"
    class="result-table"
    :items-per-page="-1"
  >
    <template v-slot:[`item.status`]="{ item }">
      <StatusChip :status="item.status" :key="item.status" />
    </template>
    <template v-slot:[`item.details`]="{ item }">
      <v-btn
        @click="details(item)"
        variant="text"
        icon="mdi-open-in-new"
        v-if="item.status === 'pass'"
      />
      <MsgTooltip :msg="item.message" v-else />
    </template>
    <template #bottom></template>
  </v-data-table>
  <v-divider />
</template>

<script setup lang="ts">
import { computed, type PropType } from 'vue'
import hash from 'object-hash'
import { useDisplay } from 'vuetify'
import { useRouter } from 'vue-router'
import type { Generation, RuleStatus } from '@/types'
import { useSessionStorage } from '@vueuse/core'
import StatusChip from './StatusChip.vue'
import MsgTooltip from './MsgTooltip.vue'

type Item = {
  id: string
  apiVersion: string
  kind: string
  resource: string
  policy: string
  rule: string
  message: string
  generatedResource: string
  status: RuleStatus
}

const props = defineProps({
  results: { type: Array as PropType<Generation[]>, default: () => [] },
  title: { type: String, default: 'Generation Results' }
})

const display = useDisplay()

const headers = computed(() => {
  if (display.mdAndUp.value) {
    return [
      { title: 'APIVersion', key: 'apiVersion', width: '15%' },
      { title: 'Kind', key: 'kind', width: '10%' },
      { title: 'Resource', key: 'resource', width: '20%' },
      { title: 'Policy', key: 'policy', width: '20%' },
      { title: 'Rule', key: 'rule', width: '25%' },
      { title: 'Status', key: 'status', width: '5%', align: 'end' },
      { title: 'Details', key: 'details', width: '5%', align: 'end' }
    ]
  }

  return [
    { title: 'Kind', key: 'kind', width: '10%' },
    { title: 'Resource', key: 'resource', width: '20%' },
    { title: 'Policy', key: 'policy', width: '30%' },
    { title: 'Rule', key: 'rule', width: '30%' },
    { title: 'Status', key: 'status', width: '10%', align: 'end' },
    { title: 'Details', key: 'details', width: '5%', align: 'end' }
  ]
})

const items = computed(() => {
  return (props.results || []).reduce<Item[]>((results, generation) => {
    const policy = generation.policy || generation.validatingAdmissionPolicy

    const rules = generation.policyResponse.rules || []

    rules.forEach((rule) => {
      const item = {
        id: '',
        apiVersion: generation.resource.apiVersion,
        kind: generation.resource.kind,
        resource: [generation.resource.metadata.namespace, generation.resource.metadata.name]
          .filter((s) => !!s)
          .join('/'),
        policy: policy.metadata.name,
        rule: rule.name,
        message: rule.message,
        generatedResource: rule.generatedResource,
        status: rule.status
      }
      item.id = hash(item)

      results.push(item)
    })

    return results
  }, [])
})

const router = useRouter()

const details = (generation: Item) => {
  useSessionStorage(`generation:${generation.id}`, generation)

  const url = router.resolve({ name: 'generation-details', params: { id: generation.id } })

  window.open(url.href, '_blank')
}
</script>

<style>
.result-table th:first-child,
.result-table td:first-child {
  padding-left: 24px !important;
}

.result-table th:last-child,
.result-table td:last-child {
  padding-right: 24px !important;
}

.v-theme--light .table-expansion {
  background-color: #eee !important;
}

.v-theme--dark .table-expansion {
  background-color: #111 !important;
}
</style>
