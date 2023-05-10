<template>
  <v-app :theme="layoutTheme">
    <v-app-bar flat border>
        <template v-slot:prepend>
          <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
        </template>
      <div class="toolbar-container">
        <div class="py-1 app-logo">
          <v-img src="/kyverno-logo.png" />
          <v-chip size="small" style="position: absolute; bottom: 14px; right: -45px;">v1.10</v-chip>
        </div>
        <h1 class="text-h4 d-none d-md-inline">Playground</h1>
      </div>
        <template v-slot:append>
          <v-btn icon="mdi-github" href="https://github.com/kyverno/playground" target="_blank" class="mr-2" title="GitHub: Kyverno Playground" />
          <PrimeButton variant="outlined" @click="drawer = !drawer" class="mr-2">Examples</PrimeButton>
          <ShareButton :policy="policy" :resource="resource" :context="context" />
          <ConfigMenu @on-reset="reset" />
        </template>
    </v-app-bar>
    <ExampleDrawer v-model="drawer" @select:example="setExample" />
    <v-main>
      <v-container fluid>
        <OnboardingAlert />
        <v-row>
          <v-col :md="7" :sm="12">
            <v-card style="height: 800px">
              <EditorToolbar title="ClusterPolicy" v-model="policy" :restore-value="loadedPolicy" />
              <PolicyEditor v-model="policy" />
            </v-card>
          </v-col>
          <v-col :md="5" :sm="12">
            <v-card style="height: 300px">
              <EditorToolbar title="Context" v-model="context" :restore-value="loadedContext" />
              <ContextEditor v-model="context" />
            </v-card>
            <v-card style="height: 487px" class="mt-3">
              <EditorToolbar title="Resource" v-model="resource" :restore-value="loadedResource">
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
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { layoutTheme } from "@/config";

import ErrorBar from "@/components/ErrorBar.vue";
import EditorToolbar from "@/components/EditorToolbar.vue";
import ExampleDrawer from "@/components/ExampleDrawer.vue";
import PolicyEditor from "@/components/PolicyEditor.vue";
import ContextEditor from "@/components/ContextEditor.vue";
import ResourceEditor from "@/components/ResourceEditor.vue";
import ValidationButton from "@/components/ValidationButton.vue";
import ValidationDialog from "@/components/ValidationDialog.vue";
import OnboardingAlert from "@/components/OnboardingAlert.vue";
import ConfigMenu from "@/components/ConfigMenu.vue";
import PrimeButton from "@/components/PrimeButton.vue";
import ShareButton from "@/components/ShareButton.vue";
import TemplateButton from "@/components/TemplateButton.vue";
import * as lzstring from "lz-string";

import { PolicyTemplate, ContextTemplate, ResourceTemplate } from "@/assets/templates";
import { EngineResponse } from '@/types'
import { onMounted } from "vue";
import { useRoute } from "vue-router";

const policy = ref<string>(PolicyTemplate);
const context = ref<string>(ContextTemplate);
const resource = ref<string>(ResourceTemplate);

const loadedContext = ref<string>('');
const loadedPolicy = ref<string>(PolicyTemplate);
const loadedResource = ref<string>(ResourceTemplate);

const setExample = (values: [string, string]) => {
  policy.value = values[0];
  resource.value = values[1];

  loadedPolicy.value = values[0]
  loadedResource.value = values[1]
};

const showResults = ref<boolean>(false);
const results = ref<EngineResponse>({ validation: [], policies: [], resources: [] });

const handleResponse = (response: EngineResponse) => {
  results.value = response
  showResults.value = true
};

const reset = () => {
  policy.value = PolicyTemplate
  context.value = ContextTemplate
  resource.value = ResourceTemplate

  loadedPolicy.value = PolicyTemplate
  loadedResource.value = ResourceTemplate

  results.value = { validation: [], policies: [], resources: [] }
  showResults.value = false

  errorText.value = ''
  showError.value = false
}

const showError = ref<boolean>(false);
const errorText = ref<string>("");

const handleError = (error: Error) => {
  errorText.value = error.message
  showError.value = true

  setTimeout(() => {
    errorText.value = ''
    showError.value = false
  }, 5000)
};

const drawer = ref<boolean>(false);

const route = useRoute()
onMounted(() => {
  const query = route.query.content as string
  if (!query) return;

  try {
    const content = JSON.parse(lzstring.decompressFromBase64(query)) as { policy: string; resource: string; context: string }

    policy.value = content.policy
    resource.value = content.resource
    context.value = content.context

    loadedPolicy.value = content.policy
    loadedResource.value = content.resource
    loadedContext.value = content.context
  } catch(err) {
    console.error('could not parse content string', err)
  }
})
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
