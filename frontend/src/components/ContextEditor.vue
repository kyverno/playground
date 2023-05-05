<template>
  <MonacoEditor
    language="json"
    theme="vs-dark"
    :value="props.modelValue"
    @update:value="(event: string) => emit('update:modelValue', event)"
    :options="options"
  />
</template>

<script setup lang="ts">
import MonacoEditor from "./MonacoEditor.vue";
import { editor } from "monaco-editor/esm/vs/editor/editor.api";

const props = defineProps({
    modelValue: { type: String, default: '' }
})

const emit = defineEmits(["update:modelValue"])

const options: editor.IStandaloneEditorConstructionOptions = {
  colorDecorators: true,
  lineHeight: 24,
  tabSize: 2,
  minimap: { enabled: false },
};
</script>
