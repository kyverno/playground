<template>
  <router-view />
</template>

<style>
.no-opacity-tooltip {
  background: rgba(var(--v-theme-surface-variant), 0.9) !important;
}

:root {
  --v-onboarding-overlay-z: 2000;
}
:root {
  --v-onboarding-step-z: 3000;
}
</style>
