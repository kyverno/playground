<template>
  <router-view />
</template>

<script setup lang="ts">
</script>

<style>
.no-opacity-tooltip {
  background: rgba(var(--v-theme-surface-variant), 0.9)!important;
}

:root {
  --v-onboarding-overlay-z: 2000;
}
:root {
  --v-onboarding-step-z: 3000;
}
</style>