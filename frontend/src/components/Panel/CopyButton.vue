<template>
  <v-tooltip
    location="bottom"
    content-class="no-opacity-tooltip"
    :text="copied ? 'copied' : 'Copy to Clipboard'"
  >
    <template v-slot:activator="{ props }">
      <v-btn @click="copy(value)" :disabled="!isSupported" v-bind="props" icon="mdi-content-copy" />
    </template>
  </v-tooltip>
</template>

<script setup lang="ts">
import { useClipboard } from '@vueuse/core'

const props = defineProps({
  value: { type: String }
})

const { copy, copied, isSupported } = useClipboard({ source: props.value })
</script>
