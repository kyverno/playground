<template>
  <v-menu location="bottom" :close-on-content-click="false">
    <template v-slot:activator="{ props }">
      <v-btn prepend-icon="mdi-cog" v-bind="props" v-if="display.mdAndUp.value">Options</v-btn>
      <v-btn icon="mdi-cog" v-bind="props" v-else />
    </template>
    <v-card min-width="250">
      <v-card-title class="text-subtitle-1">Playground Configuration</v-card-title>
      <v-divider />
      <v-card-text>
        <v-select
          :items="options.layoutThemes"
          label="Layout Theme"
          v-model="layoutTheme"
          hide-details
          variant="outlined"
          density="compact"
          item-title="name"
          item-value="theme"
          class="mt-2" />
        <v-select
          :items="options.editorThemes"
          label="Editor Theme"
          v-model="editorTheme"
          hide-details
          variant="outlined"
          density="compact"
          item-title="name"
          item-value="theme"
          class="mt-4" />
        <v-btn @click="reset" prepend-icon="mdi-delete" block class="mt-4" variant="outlined" color="error">Reset Options</v-btn>

        <v-chip v-if="version" style="width: 100%" class="mt-3 justify-center font-weight-medium" size="small" label>
          Version: {{ version.substring(0, 10) }}
        </v-chip>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { useDisplay } from 'vuetify/lib/framework.mjs'
import { useConfig } from '@/config'
import { usePreferredDark } from '@vueuse/core'

const { options, layoutTheme, editorTheme, showOnboarding } = useConfig()

const isDark = usePreferredDark()

const reset = () => {
  layoutTheme.value = isDark ? 'dark' : 'light'
  editorTheme.value = 'vs-dark'
  showOnboarding.value = true
}

const display = useDisplay()
const version = APP_VERSION
</script>
