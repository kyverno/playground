<template>
  <v-tooltip
    :model-value="clicked"
    location="bottom"
    text="Local state was resettet"
    :open-on-hover="false"
  >
    <template v-slot:activator="{ props }">
      <v-btn
        prepend-icon="mdi-delete"
        color="error"
        @click="reset"
        v-bind="props"
        >Reset</v-btn
      >
    </template>
  </v-tooltip>
</template>
<script setup lang="ts">
import { ref } from "vue";

const clicked = ref<boolean>(false);

const emit = defineEmits(['on-reset'])

const reset = () => {
  clicked.value = true
  emit('on-reset')
  setTimeout(() => {
    clicked.value = false
  }, 2000)
}
</script>
