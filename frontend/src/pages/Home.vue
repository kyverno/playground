<template>
  <v-app :theme="layoutTheme">
    <v-app-bar flat border>
      <template v-slot:prepend>
        <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      </template>
      <div class="toolbar-container">
        <router-link class="py-1 app-logo d-block" to="/">
          <v-img src="/kyverno-logo.png" />
          <v-chip size="small" style="position: absolute; bottom: 14px; right: -45px"
            >v1.10</v-chip
          >
        </router-link>
        <h1 class="text-h4 d-none d-lg-inline" style="padding-left: 200px">Playground</h1>
      </div>
      <template v-slot:append>
        <v-btn
          icon="mdi-github"
          href="https://github.com/kyverno/playground"
          target="_blank"
          class="mr-2"
          title="GitHub: Kyverno Playground"
        />
        <template v-if="display.mdAndUp.value">
          <PrimeButton variant="outlined" @click="drawer = !drawer" class="mr-2"
            >Examples</PrimeButton
          >
          <ShareButton :policy="policy" :resource="resource" :context="context" />
          <SaveButton :policy="policy" :resource="resource" :context="context" />
          <LoadButton v-model:policy="policy" v-model:resource="resource" v-model:context="context" btn-class="ml-2" />
        </template>
        <ConfigMenu />
        <MobileMenu
          v-model:policy="policy"
          v-model:resource="resource"
          v-model:context="context"
          v-if="display.smAndDown.value"
        />
      </template>
    </v-app-bar>
    <ExampleDrawer v-model="drawer" @select:example="setExample" />
    <v-main>
      <v-container fluid>
        <OnboardingAlert />
        <v-row>
          <v-col :md="7" :sm="12">
            <v-card style="height: 800px">
              <EditorToolbar
                title="ClusterPolicy"
                v-model="policy"
                :restore-value="config.policy.value"
              />
              <PolicyEditor v-model="policy" />
            </v-card>
          </v-col>
          <v-col :md="5" :sm="12">
            <v-card style="height: 300px">
              <EditorToolbar
                title="Context"
                v-model="context"
                :restore-value="config.context.value"
              />
              <ContextEditor v-model="context" />
            </v-card>
            <v-card style="height: 487px" class="mt-3">
              <EditorToolbar
                title="Resource"
                v-model="resource"
                :restore-value="config.resource.value"
              >
                <template #prepend-actions>
                  <TemplateButton @select="(template: string) => resource = template" />
                </template>
              </EditorToolbar>
              <ResourceEditor v-model="resource" />
            </v-card>
          </v-col>
        </v-row>
        <ValidationButton
          @on-response="handleResponse"
          @on-error="handleError"
          :policy="policy"
          :resource="resource"
          :context="context"
          :error-state="showError"
        />
      </v-container>
      <ErrorBar v-model="showError" :text="errorText" />
      <ValidationDialog v-model="showResults" :results="results" :policy="policy" />
      <v-card v-if="config.state.value" style="position: fixed; bottom: 0; left: 0;">
        <v-card-text class="bg-grey-darken-2">Loaded State: {{ config.state.value }}</v-card-text>
      </v-card>
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import { ref, watchEffect } from "vue";
import { useDisplay } from "vuetify/lib/framework.mjs";
import * as lzstring from "lz-string";

import { layoutTheme, useConfig } from "@/config";
import ErrorBar from "@/components/ErrorBar.vue";
import EditorToolbar from "@/components/EditorToolbar.vue";
import ExampleDrawer from "@/components/ExampleDrawer.vue";
import PolicyEditor from "@/components/PolicyEditor.vue";
import ContextEditor from "@/components/ContextEditor.vue";
import ResourceEditor from "@/components/ResourceEditor.vue";
import ValidationDialog from "@/components/ValidationDialog.vue";
import OnboardingAlert from "@/components/OnboardingAlert.vue";

import {
  TemplateButton,
  LoadButton,
  SaveButton,
  ShareButton,
  PrimeButton,
  MobileMenu,
  ConfigMenu,
  ValidationButton,
} from "@/components/buttons";

import { PolicyTemplate, ContextTemplate, ResourceTemplate } from "@/assets/templates";
import { EngineResponse } from "@/types";
import { onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();
const display = useDisplay();

const policy = ref<string>(PolicyTemplate);
const context = ref<string>(ContextTemplate);
const resource = ref<string>(ResourceTemplate);

const config = useConfig();

const setExample = (values: [string, string]) => {
  policy.value = values[0];
  resource.value = values[1];

  config.policy.value = values[0];
  config.resource.value = values[1];
  config.state.value = ''
};

const showResults = ref<boolean>(false);
const results = ref<EngineResponse>({ policies: [], resources: [] });

const handleResponse = (response: EngineResponse) => {
  results.value = response;
  showResults.value = true;
};

const showError = ref<boolean>(false);
const errorText = ref<string>("");

const handleError = (error: Error) => {
  errorText.value = error.message;
  showError.value = true;

  setTimeout(() => {
    errorText.value = "";
    showError.value = false;
  }, 5000);
};

const drawer = ref<boolean>(false);

watchEffect(() => {
  const query = route.query.content as string
  if (!query) return;

  try {
    const content = JSON.parse(lzstring.decompressFromBase64(query)) as {
      policy: string;
      resource: string;
      context: string;
    };

    policy.value = content.policy;
    resource.value = content.resource;
    context.value = content.context;

    config.policy.value = content.policy;
    config.resource.value = content.resource;
    config.context.value = content.context;
    config.state.value = ''

    router.replace({ ...route, query: {} });
  } catch (err) {
    console.error("could not parse content string", err);
  }
})

onMounted(() => {
  if (route.query.content) {
    return;
  }

  if (config.policy.value) {
    policy.value = config.policy.value;
  }
  if (config.resource.value) {
    resource.value = config.resource.value;
  }
  if (config.context.value) {
    context.value = config.context.value;
  }
});
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
  padding-left: 200px;
}

.footer {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
}
</style>
