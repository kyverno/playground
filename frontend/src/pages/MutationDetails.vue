<template>
  <v-app :theme="layoutTheme">
    <v-app-bar flat border>
      <div class="toolbar-container">
        <div class="py-1 app-logo">
          <v-img src="/kyverno-logo.png" />
          <v-chip size="small" style="position: absolute; bottom: 14px; right: -45px"
            >v1.10</v-chip
          >
        </div>
        <h1 class="text-h4 d-none d-md-inline">Playground</h1>
      </div>
      <template v-slot:append>
        <v-btn
          icon="mdi-github"
          href="https://github.com/kyverno/playground"
          target="_blank"
          class="mr-2"
          title="GitHub: Kyverno Playground"
        />
        <ConfigMenu />
        <v-btn @click="close" prepend-icon="mdi-close">Close Window</v-btn>
      </template>
    </v-app-bar>
    <v-main :class="layoutTheme === 'light' ? 'bg-grey-lighten-2' : undefined">
      <v-container fluid>
        <v-row>
          <v-col>
            <v-card title="Mutation Details">
              <v-card-text v-if="!details">
                <v-alert type="error">
                  Details not found

                  <template #append>
                    <v-btn flat color="error" :min-width="150" size="large" @click="close"
                      >Close</v-btn
                    >
                  </template>
                </v-alert>
              </v-card-text>
              <template v-else>
                <v-card-text>
                  <v-container fluid>
                    <v-row>
                      <v-col md="3"
                        ><span class="d-inline-block font-weight-bold" style="width: 70px"
                          >Policy</span
                        >
                        {{ details.policy }}</v-col
                      >
                      <v-col md="3"
                        ><span
                          class="d-inline-block font-weight-bold"
                          style="width: 120px"
                          >APIVersion/Kind</span
                        >
                        {{ details.apiVersion }}/{{ details.kind }}</v-col
                      >
                      <v-col
                        ><span class="d-inline-block font-weight-bold" style="width: 70px"
                          >Status</span
                        ><StatusChip :status="details.status"
                      /></v-col>
                    </v-row>
                    <v-row class="mt-0">
                      <v-col md="3"
                        ><span class="d-inline-block font-weight-bold" style="width: 70px"
                          >Rule</span
                        >{{ details.rule }}</v-col
                      >
                      <v-col md="3"
                        ><span
                          class="d-inline-block font-weight-bold"
                          style="width: 120px"
                          >Resource</span
                        >{{ details.resource }}</v-col
                      >
                    </v-row>
                  </v-container>
                </v-card-text>
                <DiffEditor
                  :height="600"
                  language="yaml"
                  :original="details?.originalReosurce"
                  :value="details?.patchedResource"
                  :theme="editorTheme"
                  :options="options"
                />
              </template>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import { layoutTheme, editorTheme } from "@/config";

import ConfigMenu from "@/components/ConfigMenu.vue";
import DiffEditor from "@/components/DiffEditor.vue";
import StatusChip from "@/components/StatusChip.vue";
import { useRoute } from "vue-router";
import { useSessionStorage } from "@vueuse/core";
import { RuleStatus } from "@/types";
import { ref } from "vue";
import { watch } from "vue";

const route = useRoute();

type Item = {
  apiVersion: string;
  kind: string;
  resource: string;
  policy: string;
  rule: string;
  originalReosurce: string;
  patchedResource: string;
  status: RuleStatus;
};

const details = ref<Item | undefined>(undefined);

const content = useSessionStorage<string | null>(`mutation:${route.params.id}`, null);
watch(
  content,
  (n) => {
    if (n) {
      details.value = JSON.parse(n) as Item;
    }
  },
  { immediate: true }
);

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

<style scoped>
.app-logo {
  width: 200px;
  height: 64px;
  position: absolute;
  left: 0;
}

.toolbar-container {
  width: 100%;
  max-height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  padding-left: 175px;
}

.footer {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
}
</style>