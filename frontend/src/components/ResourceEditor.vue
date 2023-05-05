<template>
  <MonacoEditor
    language="yaml"
    theme="vs-dark"
    :value="props.modelValue"
    @update:value="(event: string) => emit('update:modelValue', event)"
    :options="options"
    ref="monaco"
  />
</template>

<script setup lang="ts">
import MonacoEditor from "./MonacoEditor.vue";

const props = defineProps({
    modelValue: { type: String, default: '' }
})

const emit = defineEmits(["update:modelValue"])

const options = {
  colorDecorators: true,
  lineHeight: 24,
  tabSize: 2,
};
</script>
