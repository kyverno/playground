1<template>
  <MonacoEditor
    language="yaml"
    id="resource"
    :theme="editorTheme"
    :modelValue="props.modelValue"
    @update:modelValue="(event: string) => emit('update:modelValue', event)"
    :options="options"
    @editorDidMount="e => monaco = e"
    :uri="uri"
  />
</template>

<script setup lang="ts">
import { editor, Uri } from 'monaco-editor'
import { editorTheme } from "@/config";
import { ref, watch } from "vue";
import MonacoEditor from "./MonacoEditor.vue";

const props = defineProps({
    modelValue: { type: String, default: '' }
})

const monaco = ref<editor.ICodeEditor>();

watch(() => props.modelValue, (current, old) => {
  if (!monaco.value) return

  const currentLines = current.split(/\r\n|\r|\n/).length
  const oldLines = old.split(/\r\n|\r|\n/).length

  if (currentLines + 10 < oldLines) {
    monaco.value.setScrollPosition({scrollTop: 0});
  }
})

const uri = Uri.parse('resource.yaml')
const emit = defineEmits(["update:modelValue"])

const options = {
  colorDecorators: true,
  lineHeight: 24,
  tabSize: 2,
};
</script>
