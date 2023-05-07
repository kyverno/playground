<template>
  <v-menu location="bottom" :close-on-content-click="false">
    <template v-slot:activator="{ props }">
      <v-btn v-bind="props"> Options </v-btn>
    </template>
    <v-card min-width="250">
      <v-card-title class="text-subtitle-1"> Playground Configuration </v-card-title>
      <v-divider />
      <v-card-text>
        <v-select
          :items="config.layoutThemes"
          label="Layout Theme"
          v-model="layoutTheme"
          hide-details
          variant="outlined"
          density="compact"
          item-title="name"
          item-value="theme"
          class="mt-2"
        />
        <v-select
          :items="config.editorThemes"
          label="Editor Theme"
          v-model="editorTheme"
          hide-details
          variant="outlined"
          density="compact"
          item-title="name"
          item-value="theme"
          class="mt-4"
        />
        <v-btn block @click="emit('on-reset')" flat color="grey" class="mt-4">Reset Input</v-btn>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { editorTheme, layoutTheme, config } from "../config";

const emit = defineEmits(['on-reset'])
</script>
