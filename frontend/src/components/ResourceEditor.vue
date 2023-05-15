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
import { editor } from 'monaco-editor'
import MonacoEditor from "./MonacoEditor.vue";
import { editorTheme } from "../config";
import { Uri } from "monaco-editor";
import { watch } from "vue";
import { ref } from "vue";
import { loadedResource } from "@/composables";

const props = defineProps({
    modelValue: { type: String, default: '' }
})

const monaco = ref<typeof MonacoEditor | null>(null);

watch(loadedResource, () => {
  if (!monaco.value) return

  const edit: editor.ICodeEditor = monaco.value._getEditor();

  edit.setScrollPosition({scrollTop: 0});
})

const uri = Uri.parse('resource.yaml')
const emit = defineEmits(["update:modelValue"])

const options = {
  colorDecorators: true,
  lineHeight: 24,
  tabSize: 2,
};
</script>
