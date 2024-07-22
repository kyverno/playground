<template>
  <MonacoEditor
    language="yaml"
    id="resource"
    :theme="editorTheme"
    :modelValue="props.modelValue"
    @update:modelValue="(event: string) => emit('update:modelValue', event)"
    :options="options"
    :uri="uri"
  />
</template>

<script setup lang="ts">
import { Uri } from 'monaco-editor'
import { editorTheme } from '@/config'
import MonacoEditor from './MonacoEditor.vue'

const props = defineProps({
  modelValue: { type: String, default: '' }
})

const uri = Uri.parse('resource.yaml')
const emit = defineEmits(['update:modelValue'])

const options = {
  colorDecorators: true,
  lineHeight: 24,
  tabSize: 2
}
</script>
