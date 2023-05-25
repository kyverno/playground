1<template>
  <DiffEditor
    id="diff-policy"
    language="yaml"
    :theme="editorTheme"
    v-model="inputs.resource"
    v-model:original="inputs.oldResource"
    :options="options"
    :uri="uri"
  />
</template>

<script setup lang="ts">
import { Uri } from 'monaco-editor'
import { editorTheme } from "@/config";
import DiffEditor from "./DiffEditor.vue";
import { inputs } from '@/store';

const uri = Uri.parse('resource.yaml')

const options = {
  readOnly: false,
  originalEditable: true,
  colorDecorators: true,
  lineHeight: 24,
  tabSize: 2,
};
</script>
