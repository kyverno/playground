<template>
  <MonacoEditor
    language="yaml"
    :theme="editorTheme"
    :value="props.modelValue"
    @update:value="(event: string) => emit('update:modelValue', event)"
    :options="options"
    ref="monaco"
    :uri="uri"
  />
</template>

<script setup lang="ts">
import MonacoEditor from "./MonacoEditor.vue";
import { editorTheme } from "../config";
import { Uri } from "monaco-editor";

const props = defineProps({
    modelValue: { type: String, default: '' }
})

const uri = Uri.parse('resource.yaml')
const emit = defineEmits(["update:modelValue"])

const options = {
  colorDecorators: true,
  lineHeight: 24,
  tabSize: 2,
};
</script>
