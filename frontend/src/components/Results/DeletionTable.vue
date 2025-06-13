<template>
  <v-card-title class="my-2 text-subtitle-1">Deletion Results</v-card-title>
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
    v-model:expanded="expanded"
    :items-per-page="-1"
  >
    <template v-slot:[`item.policy`]="{ item }">
      <v-avatar size="24px" rounded="0" class="mr-2 mb-1">
        <v-img :src="`icons/${item.icon}.png`" />
      </v-avatar>
      {{ item.policy }}
    </template>
    <template v-slot:[`item.status`]="{ item }">
      <StatusChip :status="item.status" :key="item.status" />
    </template>
    <template v-slot:expanded-row="{ columns, item }">
      <tr>
        <td :colspan="columns.length" class="py-2 table-expansion">
          {{ item.message }}
        </td>
      </tr>
    </template>
    <template #bottom></template>
  </v-data-table>
  <v-divider />
</template>

<script setup lang="ts">
import { type PropType, computed, ref } from 'vue'
import hash from 'object-hash'
import { useDisplay } from 'vuetify'
import { useConfig } from '@/config'
import type { Validation, RuleStatus } from '@/types'
import StatusChip from './StatusChip.vue'

type Item = {
  id: string
  icon?: string
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

const expanded = ref<string[]>([])

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
  return (props.results || []).reduce<Item[]>((results, deletion) => {
    const policy =
      deletion.policy ||
      deletion.validatingAdmissionPolicy ||
      deletion.validatingPolicy ||
      deletion.imageValidatingPolicy ||
      deletion.deletingPolicy ||
      deletion.generatingPolicy

    if (!deletion.policyResponse.rules && !hideNoMatch.value) {
      results.push({
        id: 'id',
        apiVersion: deletion.resource.apiVersion,
        kind: deletion.resource.kind,
        resource:
          [deletion.resource?.metadata?.namespace, deletion.resource?.metadata?.name]
            .filter((s) => !!s)
            .join('/') || 'JSON Payload',
        policy: policy.metadata.name,
        rule: 'resource does not match any rule',
        message: 'no validation triggered',
        status: 'no match'
      })
      return results
    }

    const rules = deletion.policyResponse.rules || []

    rules.forEach((rule) => {
      let ruleName = rule.name
      let icon = 'kyverno'

      if (deletion.validatingAdmissionPolicy) {
        ruleName = 'N/A'
        icon = 'k8s'
      }

      if (deletion.validatingPolicy) {
        ruleName = 'N/A'
      }

      results.push({
        id: hash({ rule, resource: deletion.resource }),
        icon,
        apiVersion: deletion.resource.apiVersion,
        kind: deletion.resource.kind,
        resource:
          [deletion.resource?.metadata?.namespace, deletion.resource?.metadata?.name]
            .filter((s) => !!s)
            .join('/') || 'JSON Payload',
        policy: policy.metadata.name,
        rule: ruleName,
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
