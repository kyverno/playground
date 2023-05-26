<template>
  <v-row v-if="showOnboarding">
    <v-col>
      <v-alert type="info" variant="outlined" prominent>
        {{ options.onboarding.text }}

        <template #append>
          <v-btn flat color="info" :min-width="150" size="large" @click="() => (showOnboarding = false)">Close</v-btn>
        </template>
      </v-alert>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { useConfig } from '@/config'

const { options, showOnboarding } = useConfig()
</script>
