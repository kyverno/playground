<template>
  <v-app :theme="config.theme">
    <v-app-bar flat border>
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <div class="toolbar-container">
        <div class="py-1 app-logo">
          <v-img src="/kyverno-logo.png" />
        </div>
        <h1 style="display: inline" class="text-h4">Playground</h1>
      </div>
    </v-app-bar>
    <ExampleDrawer v-model="drawer" @select:example="setExample" />
    <v-main>
      <v-container fluid>
        <v-row>
          <v-col :md="7" :sm="12">
            <v-card style="height: 800px">
              <EditorToolbar title="ClusterPolicy" v-model="policy" />
              <PolicyEditor v-model="policy" />
            </v-card>
          </v-col>
          <v-col :md="5" :sm="12">
            <v-card style="height: 300px">
              <EditorToolbar title="Context" v-model="context" />
              <ContextEditor v-model="context" />
            </v-card>
            <v-card style="height: 487px" class="mt-3">
              <EditorToolbar title="Resource" v-model="resource" />
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
import { config } from "./config";

import ErrorBar from "./components/ErrorBar.vue";
import EditorToolbar from "./components/EditorToolbar.vue";
import ExampleDrawer from "./components/ExampleDrawer.vue";
import PolicyEditor from "./components/PolicyEditor.vue";
import ContextEditor from "./components/ContextEditor.vue";
import ResourceEditor from "./components/ResourceEditor.vue";
import ValidationButton from "./components/ValidationButton.vue";
import ValidationDialog from "./components/ValidationDialog.vue";
import { PolicyTemplate, ContextTemplate, ResourceTemplate } from "./assets/templates";
import { EngineResponse } from './types'

const policy = ref<string>(PolicyTemplate);
const context = ref<string>(ContextTemplate);
const resource = ref<string>(ResourceTemplate);

const setExample = (values: [string, string]) => {
  policy.value = values[0];
  resource.value = values[1];
};

const showResults = ref<boolean>(false);
const results = ref<EngineResponse>({ Validation: [] });

const handleResponse = (response: EngineResponse) => {
  results.value = response
  showResults.value = true
};

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
}
</style>
