<template>
  <v-app-bar flat border>
    <template v-slot:prepend>
      <slot name="prepend-actions" />
    </template>
    <div class="toolbar-container">
      <router-link class="py-1 app-logo d-block" to="/">
        <template v-if="display.mdAndUp.value">
          <v-img src="/kyverno-logo.png" />
          <v-chip size="small" style="position: absolute; bottom: 14px; right: -45px">v1.10</v-chip>
        </template>
        <template v-if="display.smAndDown.value">
          <v-img src="/favicon.png" width="80" />
          <v-chip size="small" style="position: absolute; bottom: 16px; left: 80px">Kyverno v1.10</v-chip>
        </template>
      </router-link>
      <h1 class="text-h4 d-none d-lg-inline" style="padding-left: 200px">Kyverno Playground</h1>
    </div>
    <template v-slot:append>
      <v-btn
        href="http://kyverno.io"
        target="_blank"
        class="mr-1"
        title="Kyverno Documentation"
        v-if="display.mdAndUp.value"
      >
        <v-img src="/favicon.png" :height="24" :width="24" alt="Kyverno Logo" class="mr-2" /> Docs
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
        <slot name="append-actions" />
      <slot name="mobile-actions" v-if="display.smAndDown.value" />
    </template>
  </v-app-bar>
</template>

<script setup lang="ts">
import { useDisplay } from "vuetify/lib/framework.mjs";
import ConfigMenu from "./ConfigMenu.vue";

const display = useDisplay();
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
  padding-left: 200px;
}
</style>
