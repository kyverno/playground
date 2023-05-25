<template>
  <v-dialog v-model="dialog">
    <template v-slot:activator="{ props }">
      <v-btn flat color="primary" block v-bind="props">
        Custom Resource Definitions
      </v-btn>
    </template>
    <v-card title="Custom Resource Definitions">
      <v-card-text class="px-0 py-0" >
        <MonacoEditor :theme="editorTheme" :value="inputs.customResourceDefinitions"
          @update:value="(event: string) => inputs.customResourceDefinitions = event" language="yaml" />
      </v-card-text>
      <v-card-actions>
        <v-btn @click="dialog = false">Close</v-btn>
        <v-spacer />
        <v-btn @click="dialog = false">Ok</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { editorTheme } from "@/config"
import { inputs } from "@/store";
import { MonacoEditor } from '.';

const dialog = ref<boolean>(false);
</script>
