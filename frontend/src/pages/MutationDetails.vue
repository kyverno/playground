<template>
  <v-app :theme="layoutTheme">
    <AppBar>
      <template #append-actions>
        <v-btn @click="close" prepend-icon="mdi-close" class="d-md-flex d-none"
          >Close Window</v-btn
        >
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
                    <v-btn
                      flat
                      color="error"
                      :min-width="150"
                      size="large"
                      @click="close"
                      >Close</v-btn
                    >
                  </template>
                </v-alert>
              </v-card-text>
              <template v-else>
                <v-card-text>
                  <v-container fluid>
                    <v-row>
                      <v-col md="3" cols="6"
                        ><span
                          class="d-inline-block font-weight-bold"
                          style="width: 50px"
                          >Policy</span
                        >
                        {{ details.policy }}</v-col
                      >
                      <v-col md="3" cols="6"
                        ><span
                          class="d-inline-block font-weight-bold"
                          style="width: 120px"
                          >APIVersion/Kind</span
                        >
                        {{ details.apiVersion }}/{{ details.kind }}</v-col
                      >
                      <v-col cols="3" class="d-none d-md-flex align-center"
                        ><span
                          class="d-inline-block font-weight-bold"
                          style="width: 50px"
                          >Status</span
                        >
                        <StatusChip :status="details.status" />
                      </v-col>
                    </v-row>
                    <v-row class="mt-0">
                      <v-col md="3" cols="6"
                        ><span
                          class="d-inline-block font-weight-bold"
                          style="width: 50px"
                          >Rule</span
                        >{{ details.rule }}</v-col
                      >
                      <v-col md="3" cols="6"
                        ><span
                          class="d-inline-block font-weight-bold"
                          style="width: 120px"
                          >Resource</span
                        >{{ details.resource }}</v-col
                      >
                    </v-row>
                    <v-row>
                      <v-col>
                        <span
                          class="d-inline-block font-weight-bold"
                          style="width: 50px"
                          >Status</span
                        >
                        <StatusChip :status="details.status" />
                      </v-col>
                    </v-row>
                  </v-container>
                </v-card-text>
                <v-card-text v-if="details.status !== 'pass'">{{
                  details.message
                }}</v-card-text>
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
import { ref, watch } from "vue";
import { useRoute } from "vue-router";
import { useSessionStorage } from "@vueuse/core";
import { RuleStatus } from "@/types";
import { layoutTheme, editorTheme } from "@/config";

import DiffEditor from "@/components/DiffEditor.vue";
import StatusChip from "@/components/StatusChip.vue";
import AppBar from "@/components/AppBar.vue";

const route = useRoute();

type Item = {
  apiVersion: string;
  kind: string;
  resource: string;
  policy: string;
  rule: string;
  message: string;
  originalReosurce: string;
  patchedResource: string;
  status: RuleStatus;
};

const details = ref<Item | undefined>(undefined);

const content = useSessionStorage<string | null>(
  `mutation:${route.params.id}`,
  null
);
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
