<template>
  <MonacoEditor
    id="context"
    language="yaml"
    :theme="editorTheme"
    :modelValue="props.modelValue"
    @update:modelValue="(event: string) => emit('update:modelValue', event)"
    :options="options"
    ref="monaco"
    :uri="uri"
  />
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { editor } from "monaco-editor/esm/vs/editor/editor.api";
import { Uri } from "monaco-editor";
import { editorTheme } from "@/config";
import { loadedContext } from "@/composables";
import MonacoEditor from "./MonacoEditor.vue";

const props = defineProps({
    modelValue: { type: String, default: '' }
})

const monaco = ref<typeof MonacoEditor | null>(null);

watch(loadedContext, () => {
  if (!monaco.value) return

  const edit: editor.ICodeEditor = monaco.value._getEditor();

  edit.setScrollPosition({scrollTop: 0});
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
