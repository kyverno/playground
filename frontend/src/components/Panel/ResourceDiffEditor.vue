1<template>
  <DiffEditor
    language="yaml"
    :theme="editorTheme"
    :value="props.modelValue"
    :original="inputs.oldResource"
    @update:value="(event: string) => inputs.resource = event"
    @update:original="(event: string) => inputs.oldResource = event"
    :options="options"
    ref="monaco"
    :uri="uri"
  />
</template>

<script setup lang="ts">
import { editor, Uri } from 'monaco-editor'
import { editorTheme } from "@/config";
import { ref, watch } from "vue";
import DiffEditor from "../Details/DiffEditor.vue";
import { inputs } from '@/store';

const props = defineProps({
    modelValue: { type: String, default: '' }
})

const monaco = ref<typeof DiffEditor | null>(null);

watch(() => props.modelValue, (current, old) => {
  if (!monaco.value) return

  const edit: editor.ICodeEditor = monaco.value._getEditor();
  if (!edit) return

  const currentLines = current.split(/\r\n|\r|\n/).length
  const oldLines = old.split(/\r\n|\r|\n/).length

  if (currentLines + 10 < oldLines) {
    edit.setScrollPosition({scrollTop: 0});
  }
})

const uri = Uri.parse('resource.yaml')
// const emit = defineEmits(["update:modelValue"])

const options = {
  readOnly: false,
  originalEditable: true,
  colorDecorators: true,
  lineHeight: 24,
  tabSize: 2,
};
</script>
