<template>
  <v-card-title class="my-2 text-subtitle-1">Validation Results</v-card-title>
  <v-divider />
  <!-- @vue-ignore -->
  <v-data-table
    density="default"
    hover
    :items="items"
    item-value="id"
    :headers="headers as any"
    class="result-table"
    show-expand
    v-model:expanded="(expanded as any)"
    :items-per-page="-1">
    <template v-slot:[`item.status`]="{ item }">
      <StatusChip :status="item.raw.status" :key="item.raw.status" />
    </template>
    <template v-slot:expanded-row="{ columns, item }">
      <tr>
        <td :colspan="columns.length" class="py-2 table-expansion">
          {{ item.raw.message }}
        </td>
      </tr>
    </template>
    <template #bottom></template>
  </v-data-table>
  <v-divider />
</template>

<script setup lang="ts">
import { PropType, computed, ref } from 'vue'
import hash from 'object-hash'
import { useDisplay } from 'vuetify/lib/framework.mjs'
import { useConfig } from '@/config'
import { Validation, RuleStatus } from '@/types'
import StatusChip from './StatusChip.vue'

type Item = {
  id: string;
  apiVersion: string
  kind: string
  resource: string
  policy: string
  rule: string
  message: string
  status: RuleStatus
}

const props = defineProps({
  results: { type: Array as PropType<Validation[]>, default: () => [] }
})

const expanded = ref<number[]>([])

const display = useDisplay()

const headers = computed(() => {
  if (display.mdAndUp.value) {
    return [
      { title: 'APIVersion', key: 'apiVersion', width: '15%' },
      { title: 'Kind', key: 'kind', width: '10%' },
      { title: 'Resource', key: 'resource', width: '20%' },
      { title: 'Policy', key: 'policy', width: '25%' },
      { title: 'Rule', key: 'rule', width: '25%' },
      { title: 'Status', key: 'status', width: '5%', align: 'end' }
    ]
  }

  return [
    { title: 'Kind', key: 'kind', width: '10%' },
    { title: 'Resource', key: 'resource', width: '20%' },
    { title: 'Policy', key: 'policy', width: '30%' },
    { title: 'Rule', key: 'rule', width: '30%' },
    { title: 'Status', key: 'status', width: '10%', align: 'end' }
  ]
})

const { hideNoMatch } = useConfig()

const items = computed(() => {
  return (props.results || []).reduce<Item[]>((results, validation) => {
    if (!validation.policyResponse.rules && !hideNoMatch.value) {
      results.push({
        id: 'id',
        apiVersion: validation.resource.apiVersion,
        kind: validation.resource.kind,
        resource: [validation.resource.metadata.namespace, validation.resource.metadata.name].filter((s) => !!s).join('/'),
        policy: validation.policy.metadata.name,
        rule: 'resource does not match any rule',
        message: 'no validation triggered',
        status: 'no match'
      })
      return results
    }

    const rules = validation.policyResponse.rules || []

    rules.forEach((rule) => {
      results.push({
        id: hash({ rule, resource: validation.resource }),
        apiVersion: validation.resource.apiVersion,
        kind: validation.resource.kind,
        resource: [validation.resource.metadata.namespace, validation.resource.metadata.name].filter((s) => !!s).join('/'),
        policy: validation.policy.metadata.name,
        rule: rule.name,
        message: rule.message,
        status: rule.status
      })
    })

    return results
  }, [])
})
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
