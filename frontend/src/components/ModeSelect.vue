<template>
  <v-select
    clearable
    variant="outlined"
    density="comfortable"
    hide-details
    :items="modes"
    @update:modelValue="(ns: string) => emit('update:modelValue', ns)"
    :model-value="prop.modelValue"
    label="Mode"
  />
</template>

<script setup lang="ts">
const prop = defineProps({
  modelValue: { type: String },
})

const emit = defineEmits(['update:modelValue'])

const modes = ['replace', 'append']
</script>
