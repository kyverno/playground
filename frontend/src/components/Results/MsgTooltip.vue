<template>
  <v-tooltip :text="msg" content-class="no-opacity-tooltip">
    <template v-slot:activator="{ props }">
      <v-btn v-bind="props" icon="mdi-alert-circle-outline" variant="text" small />
    </template>
  </v-tooltip>
</template>

<script setup lang="ts">
defineProps({
  msg: { type: String, required: true }
})
</script>

<style>
.bg-tooltip {
  background: rgba(var(--v-theme-surface-variant), 0.9) !important;
}
</style>
