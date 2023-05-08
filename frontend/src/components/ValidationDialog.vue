<template>
  <v-dialog
    :model-value="props.modelValue"
    width="90%"
    @update:model-value="emit('update:modelValue', false)"
  >
    <v-card>
      <v-toolbar color="transparent">
        <v-toolbar-title>Validation Results</v-toolbar-title>
        <template v-slot:append>
          <v-btn flat icon="mdi-close" @click="emit('update:modelValue', false)"></v-btn>
        </template>
      </v-toolbar>
      <v-divider />
      <v-data-table
        density="default"
        hover
        :items="items"
        :headers="headers"
        class="result-table"
        show-expand
        v-model:expanded="expanded"
        :items-per-page="-1"
      >
        <template v-slot:[`item.status`]="{ item }">
          <StatusChip :status="item.raw.status" />
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
      <v-card-actions>
        <v-btn color="error" @click="emit('update:modelValue', false)">Close</v-btn>
        <v-spacer />
        <DownloadBtn
          variant="tonal"
          :content="policy"
          :filename="filename"
          label="Download Policy"
        />
        <v-tooltip
          :model-value="copied"
          location="top"
          text="Copied"
          :open-on-hover="false"
        >
          <template v-slot:activator="{ props }">
            <v-btn
              variant="tonal"
              color="secondary"
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
import { EngineResponse, RuleStatus } from "../types";
import StatusChip from "./StatusChip.vue";
import DownloadBtn from "./DownloadBtn.vue";
import { PropType } from "vue";
import { useClipboard } from "@vueuse/core";
import { computed, ref } from "vue";
import { useDisplay } from "vuetify/lib/framework.mjs";

type Item = {
  apiVersion: string;
  kind: string;
  resource: string;
  policy: string;
  rule: string;
  message: string;
  status: RuleStatus;
};

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  results: { type: Object as PropType<EngineResponse>, required: true },
  policy: { type: String, default: "" },
});

const expanded = ref<number[]>([]);

const display = useDisplay();

const headers = computed(() => {
  if (display.mdAndUp.value) {
    return [
      { title: "APIVersion", key: "apiVersion", width: "15%" },
      { title: "Kind", key: "kind", width: "10%" },
      { title: "Resource", key: "resource", width: "20%" },
      { title: "Policy", key: "policy", width: "25%" },
      { title: "Rule", key: "rule", width: "25%" },
      { title: "Status", key: "status", width: "5%", align: "end" },
    ];
  }

  return [
    { title: "Kind", key: "kind", width: "10%" },
    { title: "Resource", key: "resource", width: "20%" },
    { title: "Policy", key: "policy", width: "30%" },
    { title: "Rule", key: "rule", width: "30%" },
    { title: "Status", key: "status", width: "10%", align: "end" },
  ];
});

const filename = computed(
  () => `${(props.results.Validation || [{}])[0].policy?.metadata.name || "policy"}.yaml`
);

const items = computed(() => {
  return props.results.Validation.reduce<Item[]>((results, validation) => {
    (validation.policyResponse.rules || []).forEach((rule) => {
      results.push({
        apiVersion: validation.resource.apiVersion,
        kind: validation.resource.kind,
        resource: [
          validation.resource.metadata.namespace,
          validation.resource.metadata.name,
        ]
          .filter((s) => !!s)
          .join("/"),
        policy: validation.policy.metadata.name,
        rule: rule.name,
        message: rule.message,
        status: rule.status,
      });
    });

    return results;
  }, []);
});

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
