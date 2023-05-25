<template>
  <v-app :theme="layoutTheme">
    <AppBar>
      <template #append-actions>
        <v-btn @click="close" prepend-icon="mdi-close" class="d-md-flex d-none">Close Window</v-btn>
      </template>
      <template #mobile-actions>
        <v-btn @click="close" icon="mdi-close"></v-btn>
      </template>
    </AppBar>
    <v-main :class="layoutTheme === 'light' ? 'bg-grey-lighten-2' : undefined">
      <v-container fluid>
        <v-row>
          <v-col>
            <v-card title="Details">
              <v-card-text v-if="!details">
                <v-alert type="error">
                  Details not found
                  <template #append>
                    <v-btn flat color="error" :min-width="150" size="large" @click="close">Close</v-btn>
                  </template>
                </v-alert>
              </v-card-text>
              <template v-else>
                <v-card-text>
                  <RuleDetails :details="details" />
                </v-card-text>
                <v-card-text v-if="details.status !== 'pass'">{{ details.message }}</v-card-text>
                <ManocoEditor :height="600" language="yaml" :value="details.generatedResource" :theme="editorTheme" :options="options" id="generation" />
              </template>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { useRoute } from "vue-router";
import { useSessionStorage } from "@vueuse/core";
import { layoutTheme, editorTheme } from "@/config";
import { RuleStatus } from "@/types";
import { AppBar } from "@/components/AppBar";
import { RuleDetails } from "@/components/Details";
import ManocoEditor from "@/components/Panel/MonacoEditor.vue";

const route = useRoute();

type Item = {
  apiVersion: string;
  kind: string;
  resource: string;
  policy: string;
  rule: string;
  message: string;
  generatedResource: string;
  status: RuleStatus;
};

const details = ref<Item | undefined>(undefined);

const content = useSessionStorage<string | null>(`generation:${route.params.id}`, null);
watch(content, (n) => {
  if (!n) return;
  details.value = JSON.parse(n) as Item;
}, { immediate: true });

const close = () => {
  content.value = null;
  window.close();
};

const options = {
  readOnly: true,
  colorDecorators: true,
  lineHeight: 24,
  tabSize: 2,
};
</script>
