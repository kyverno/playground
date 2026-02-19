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
import type { Validation, RuleStatus, Resource } from '@/types'
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
  authzStatus?: number
  authzMessage?: string
}

const props = defineProps({
  results: { type: Array as PropType<Validation[]>, default: () => [] }
})

const expanded = ref<string[]>([])

const { hideNoMatch } = useConfig()

const resource = (mode: string, resource: Resource): string => {
  if (mode === 'Envoy') {
    return 'authv3.CheckRequest'
  } else if (mode === 'JSON') {
    return 'JSON Payload'
  } else if (mode === 'HTTP') {
    return 'HTTP Request'
  }

  return [resource?.metadata?.namespace, resource?.metadata?.name].filter((s) => !!s).join('/')
}

const items = computed(() => {
  return (props.results || []).reduce<Item[]>((results, validation) => {
    if (!validation.policyResponse.rules && !hideNoMatch.value) {
      results.push({
        id: 'id',
        apiVersion: validation.resource.apiVersion,
        kind: validation.resource.kind,
        resource: resource(validation.policy.mode, validation.resource),
        policy: validation.policy.name,
        rule: 'resource does not match any rule',
        message: 'no validation triggered',
        status: 'no match'
      })
      return results
    }

    const rules = validation.policyResponse.rules || []

    rules.forEach((rule) => {
      let ruleName = rule.name
      let icon = 'kyverno'

      if (validation.policy.kind === 'ValidatingAdmissionPolicy') {
        ruleName = 'N/A'
        icon = 'k8s'
      }

      if (validation.policy.kind.endsWith('ValidatingPolicy')) {
        ruleName = 'N/A'
      }

      results.push({
        id: hash({ rule, resource: validation.resource }),
        icon,
        apiVersion: validation.resource.apiVersion,
        kind: validation.resource.kind,
        resource: resource(validation.policy.mode, validation.resource),
        policy: validation.policy.name,
        rule: ruleName,
        message: rule.message,
        status: rule.status,
        authzStatus: rule.responseStatus?.code,
        authzMessage: rule.responseStatus?.message
      })
    })

    return results
  }, [])
})

const display = useDisplay()

const headers = computed(() => {
  const mode = props.results?.[0]?.policy.mode

  if ('Envoy' === mode) {
    if (display.mdAndUp.value) {
      return [
        { title: 'Resource', key: 'resource', width: '20%' },
        { title: 'Policy', key: 'policy', width: '30%' },
        { title: 'Response Code', key: 'authzStatus', width: '20%' },
        { title: 'Response Message', key: 'authzMessage', width: '20%' },
        { title: 'Status', key: 'status', width: '10%', align: 'end' }
      ]
    }

    return [
      { title: 'Policy', key: 'policy', width: '70%' },
      { title: 'Status', key: 'status', width: '30%', align: 'end' }
    ]
  }

  if (['HTTP', 'JSON'].includes(mode)) {
    if (display.mdAndUp.value) {
      return [
        { title: 'Resource', key: 'resource', width: '20%' },
        { title: 'Policy', key: 'policy', width: '60%' },
        { title: 'Status', key: 'status', width: '20%', align: 'end' }
      ]
    }
    return [
      { title: 'Policy', key: 'policy', width: '70%' },
      { title: 'Status', key: 'status', width: '30%', align: 'end' }
    ]
  }

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
