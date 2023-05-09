<template>
  <MonacoEditor
    language="yaml"
    :theme="editorTheme"
    :value="props.modelValue"
    @update:value="(event: string) => emit('update:modelValue', event)"
    :options="options"
    :uri="uri"
  />
</template>

<script setup lang="ts">
import MonacoEditor from "./MonacoEditor.vue";
import { editor } from "monaco-editor/esm/vs/editor/editor.api";
import { Uri } from "monaco-editor";
import { editorTheme } from "../config";

const props = defineProps({
    modelValue: { type: String, default: '' }
})

const emit = defineEmits(["update:modelValue"])
const uri = Uri.parse("context.yaml");

const options: editor.IStandaloneEditorConstructionOptions = {
  colorDecorators: true,
  lineHeight: 24,
  tabSize: 2,
  minimap: { enabled: false },
};
</script>
