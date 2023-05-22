<template>
  <v-dialog
    :model-value="props.modelValue"
    width="90%"
    @update:model-value="emit('update:modelValue', false)"
  >
    <v-card>
      <v-toolbar color="transparent">
        <v-toolbar-title>Results</v-toolbar-title>
        <template v-slot:append>
          <v-checkbox variant="compact" label="Hide no match results" hide-details v-model="hideNoMatch" class="mr-4" />
          <v-btn flat icon="mdi-close" @click="emit('update:modelValue', false)"></v-btn>
        </template>
      </v-toolbar>
      <v-divider />
        <v-card-text v-if="!hasResults">
          <v-alert type="warning" variant="outlined">
            No resource matched any rule of the provided policies. Please check your manifests.
          </v-alert>
        </v-card-text>
        <MutationTable :results="mutations" v-if="hasResults && mutations.length" />
        <MutationTable :results="verifications" v-if="hasResults && verifications.length" title="ImageVerification Results" />
        <ValidationTable :results="validations" v-if="hasResults && validations.length" />
        <GenerationTable :results="generations" v-if="hasResults && generations.length" />
      <v-card-actions>
        <v-btn color="error" @click="emit('update:modelValue', false)">Close</v-btn>
        <v-spacer />
        <v-tooltip
          :model-value="copied"
          location="top"
          text="Copied"
          :open-on-hover="false"
        >
          <template v-slot:activator="{ props }">
            <v-btn
              variant="tonal"
              :color="btnColor"
              @click="copy(policy)"
              :disabled="!isSupported"
              v-bind="props"
              >Copy Policy to Clipboard</v-btn
            >
          </template>
        </v-tooltip>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { computed, PropType } from "vue";
import { EngineResponse } from "@/types";
import { useClipboard } from "@vueuse/core";
import { useConfig, btnColor } from "@/config";
import ValidationTable from './ValidationTable.vue'
import MutationTable from './MutationTable.vue'
import GenerationTable from './GenerationTable.vue'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  results: { type: Object as PropType<EngineResponse>, required: true },
  policy: { type: String, default: "" },
});

const hasResults = computed(() => {
  return (props.results.validation || []).some((v) => v.policyResponse.rules !== null && v.policyResponse.rules.length > 0) || 
    (props.results.mutation || []).some((v) => v.policyResponse.rules !== null && v.policyResponse.rules.length > 0) || 
    (props.results.imageVerification || []).some((v) => v.policyResponse.rules !== null && v.policyResponse.rules.length > 0) || 
    (props.results.generation || []).some((v) => v.policyResponse.rules !== null && v.policyResponse.rules.length > 0)
})

const validations = computed(() => {
  if (hasResults.value) {
    return (props.results.validation || []).filter(v => v.policyResponse.rules)
  }

  return (props.results.validation || [])
})

const mutations = computed(() => {
  if (hasResults.value) {
    return (props.results.mutation || []).filter(v => v.policyResponse.rules)
  }

  return (props.results.mutation || [])
})

const verifications = computed(() => {
  if (hasResults.value) {
    return (props.results.imageVerification || []).filter(v => v.policyResponse.rules)
  }

  return (props.results.imageVerification || [])
})

const generations = computed(() => {
  if (hasResults.value) {
    return (props.results.generation || []).filter(v => v.policyResponse.rules)
  }

  return (props.results.generation || [])
})

const { hideNoMatch } = useConfig()

const { copy, copied, isSupported } = useClipboard({ source: props.policy });

const emit = defineEmits(["update:modelValue"]);
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
  background-color: #eee!important;
}

.v-theme--dark .table-expansion {
  background-color: #111!important;
}
</style>
