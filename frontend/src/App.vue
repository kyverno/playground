<template>
  <v-app :theme="config.theme">
    <v-toolbar flat color="white" border>
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <div class="toolbar-container">
        <div class="py-1 app-logo">
          <v-img src="/kyverno-logo.png" />
        </div>
        <h1 style="display: inline" class="text-h4">Playground</h1>
      </div>
    </v-toolbar>
    <ExampleDrawer v-model="drawer" @select:example="setExample" />
    <v-container fluid class="mt-1">
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
      />
    </v-container>
  </v-app>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { config } from "./config";

import EditorToolbar from "./components/EditorToolbar.vue";
import ExampleDrawer from "./components/ExampleDrawer.vue";
import PolicyEditor from "./components/PolicyEditor.vue";
import ContextEditor from "./components/ContextEditor.vue";
import ResourceEditor from "./components/ResourceEditor.vue";
import ValidationButton from "./components/ValidationButton.vue";
import { PolicyTemplate, ContextTemplate, ResourceTemplate } from "./assets/templates";

const policy = ref<string>(PolicyTemplate);
const context = ref<string>(ContextTemplate);
const resource = ref<string>(ResourceTemplate);

const setExample = (values: [string, string]) => {
  policy.value = values[0];
  resource.value = values[1];
};

const handleResponse = (response: Object) => {
  console.log(response);
};

const handleError = (error: Error) => {
  console.error(error);
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
