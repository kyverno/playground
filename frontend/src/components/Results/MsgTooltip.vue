<template>
  <v-tooltip :text="msg" content-class="no-opacity-tooltip">
    <template v-slot:activator="{ props }">
      <v-btn v-bind="props" @click="copy(msg)" :icon="copied ? 'mdi-alert-circle-check-outline' : 'mdi-alert-circle-outline'" variant="text" small />
    </template>
  </v-tooltip>
</template>

<script setup lang="ts">
import { useClipboard } from '@vueuse/core';

const props = defineProps({
  msg: { type: String, required: true },
})

const { copy, copied } = useClipboard({ source: props.msg })
</script>

<style>
.bg-tooltip {
  background: rgba(var(--v-theme-surface-variant), 0.9) !important;
}
</style>
