<template>
    <v-snackbar :model-value="props.modelValue" @update:model-value="(event: boolean) => emit('update:modelValue', event)" color="error">
      <span v-html="text" />
      <template v-slot:actions>
        <v-btn
          variant="text"
          @click="emit('update:modelValue', false)"
        >
          Close
        </v-btn>
      </template>
    </v-snackbar>
</template>

<script setup lang="ts">
const props = defineProps({
    modelValue: { type: Boolean, default: false },
    text: { type: String, default: '' }
})

const emit = defineEmits(["update:modelValue"])
</script>