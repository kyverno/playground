<template>
  <v-row v-if="open">
    <v-col>
      <v-alert type="info" variant="outlined" prominent v-model="open">
        Welcome to Kyverno Playground, please note that in the current version only validation
        policies are supported

        <template #append>
          <v-btn flat color="info" :min-width="150" size="large"  @click="() => open = false">Close</v-btn>
        </template>
      </v-alert>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { useLocalStorage } from "@vueuse/core";

const open = useLocalStorage("onboarding:open", true);
</script>
