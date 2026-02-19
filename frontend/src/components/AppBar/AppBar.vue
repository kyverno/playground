<template>
  <v-app-bar flat border>
    <template v-slot:prepend>
      <slot name="prepend-actions" />

      <div class="ml-2 my-1 position-relative">
        <template v-if="display.mdAndUp.value">
          <img src="/kyverno-logo.png" class="logo" />

          <v-chip
            v-if="!config.versions?.length"
            size="small"
            style="position: absolute; bottom: 20px; right: -90px"
          >
            v1.17.1
          </v-chip>
        </template>
        <template v-if="display.smAndDown.value">
          <img src="/favicon.png" class="logo" />
          <v-chip
            v-if="!config.versions?.length"
            size="small"
            style="position: absolute; bottom: 16px; left: 90px"
          >
            Kyverno v1.17.1
          </v-chip>
        </template>
      </div>
    </template>

    <div>
      <v-menu v-if="config.versions.length" open-on-hover>
        <template v-slot:activator="{ props }">
          <v-btn variant="outlined" class="text-none" rounded="xl" v-bind="props">v1.17.1</v-btn>
        </template>

        <v-list variant="flat" class="my-0 py-0 border">
          <template v-for="(version, i) in config.versions" :key="version.name">
            <v-list-item :href="version.url">
              <v-list-item-title>{{ version.name }}</v-list-item-title>
            </v-list-item>
            <v-divider v-if="i < config.versions.length - 1" />
          </template>
        </v-list>
      </v-menu>
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
import { useAPIConfig } from '@/composables/api'

const display = useDisplay()
const { config } = useAPIConfig()
</script>

<style>
.logo {
  max-height: 58px;
  margin-top: 6px;
}

.v-toolbar__prepend {
  margin-right: 24px;
}
</style>
