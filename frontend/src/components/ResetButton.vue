<template>
  <v-tooltip
    :model-value="clicked"
    location="bottom"
    text="local state resettet"
    :open-on-hover="false"
    content-class="no-opacity-tooltip"
  >
    <template v-slot:activator="{ props }">
      <v-btn
        :variant="variant"
        :block="block"
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
import { PropType, ref } from "vue";

const clicked = ref<boolean>(false);

const emit = defineEmits(['on-reset'])

defineProps({
  variant: { type: String as PropType<"outlined" | "text"> },
  block: { type: Boolean }
})

const reset = () => {
  clicked.value = true
  emit('on-reset')
  setTimeout(() => {
    clicked.value = false
  }, 2000)
}
</script>
