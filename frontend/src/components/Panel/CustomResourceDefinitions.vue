<template>
  <v-dialog v-model="dialog">
    <template v-slot:activator="{ props }">
      <v-btn flat color="primary" block v-bind="props">Custom Resource Definitions</v-btn>
    </template>
    <v-card title="Custom Resource Definitions">
      <MonacoEditor id="crds" language="yaml" :theme="editorTheme" v-model="content" :height="200" :options="options" />
      <v-card-actions>
        <v-btn @click="close">Close</v-btn>
        <v-spacer />
        <v-btn @click="save">Save</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { editorTheme } from '@/config'
import { inputs } from '@/store'
import { MonacoEditor } from '.'
import { editor } from 'monaco-editor/esm/vs/editor/editor.api'

const dialog = ref<boolean>(false)
const content = ref<string>(inputs.customResourceDefinitions)

watch(dialog, (open) => {
  if (!open) return

  content.value = inputs.customResourceDefinitions
})

const options: editor.IStandaloneEditorConstructionOptions = {
  colorDecorators: true,
  lineHeight: 24,
  tabSize: 2,
  minimap: { enabled: false }
}

const close = () => {
  dialog.value = false
}

const save = () => {
  inputs.customResourceDefinitions = content.value
  close()
}
</script>
