<template>
  <v-app-bar flat border>
    <template v-slot:prepend>
      <slot name="prepend-actions" />
    </template>
    <div class="app-logo">
      <template v-if="display.mdAndUp.value">
        <v-img src="/kyverno-logo.png" />
        <v-chip size="small" style="position: absolute; bottom: 12px; right: -80px">v1.12.6</v-chip>
      </template>
      <template v-if="display.smAndDown.value">
        <v-img src="/favicon.png" width="80" />
        <v-chip size="small" style="position: absolute; bottom: 16px; left: 80px"
          >Kyverno v1.12.5</v-chip
        >
      </template>
    </div>
    <template v-slot:append>
      <v-btn
        href="http://kyverno.io"
        target="_blank"
        class="mr-1"
        title="Kyverno Documentation"
        v-if="display.mdAndUp.value"
      >
        <v-img src="/favicon.png" :height="24" :width="24" alt="Kyverno Logo" class="mr-2" />
        Docs
      </v-btn>
      <v-btn
        icon="mdi-github"
        href="https://github.com/kyverno/playground"
        target="_blank"
        class="mr-2"
        title="GitHub: Kyverno Playground"
      />
      <template v-if="display.mdAndUp.value">
        <slot name="desktop-actions" />
      </template>
      <ConfigMenu />
      <template v-if="display.mdAndUp.value">
        <slot name="desktop-append" />
      </template>
      <slot name="append-actions" />
      <slot name="mobile-actions" v-if="display.smAndDown.value" />
    </template>
  </v-app-bar>
</template>

<script setup lang="ts">
import { useDisplay } from 'vuetify'
import ConfigMenu from './ConfigMenu.vue'

const display = useDisplay()
</script>

<style scoped>
.app-logo {
  width: 200px;
  height: 62px;
  position: absolute;
  left: 70px;
}
</style>
