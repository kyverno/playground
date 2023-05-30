<template>
  <v-dialog v-model="dialog" height="90vh">
    <template v-slot:activator="{ props }">
      <v-btn flat color="primary" block v-bind="props" :id="`${id}-btn`">{{ title }}</v-btn>
    </template>
    <v-card>
      <v-toolbar color="transparent">
        <v-toolbar-title>
          {{ title }}
          <v-tooltip :text="info" content-class="no-opacity-tooltip" v-if="info">
            <template v-slot:activator="{ props }">
              <v-btn v-bind="props" icon="mdi-information-outline" variant="text" size="small" />
            </template>
          </v-tooltip>
        </v-toolbar-title>
        <template v-slot:append>
          <slot name="append" :update="update" :content="content" />
          <v-btn flat icon="mdi-close" @click="dialog = false"></v-btn>
        </template>
      </v-toolbar>
      <v-divider />
      <v-card-text>
        <MonacoEditor
          :id="id"
          :language="language"
          :theme="editorTheme"
          v-model="content"
          height="calc(90vh - 192px)"
          :options="options"
          class="border"
          @editor-did-mount="initRefresh"
          :uri="uriParsed" />
      </v-card-text>
      <v-card-actions>
        <v-btn @click="dialog = false">Close</v-btn>
        <v-spacer />
        <v-btn @click="() => update('')" prepend-icon="mdi-delete">Clear</v-btn>
        <slot name="actions" :update="update" :content="content" />
        <v-spacer />
        <LegendMenu />
        <v-btn @click="save">Save</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
  <v-snackbar close-on-content-click color="success" v-model="success" icon="mdi-content-save">
    <div class="text-center">{{ title }} updated</div>
  </v-snackbar>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { editorTheme } from '@/config'
import MonacoEditor from '@/components/Panel/MonacoEditor.vue'
import { KeyCode, KeyMod, Uri, editor } from 'monaco-editor/esm/vs/editor/editor.api'
import LegendMenu from '../LegendMenu.vue'

const props = defineProps({
  id: { type: String, required: true },
  title: { type: String, required: true },
  info: { type: String },
  language: { type: String, default: 'yaml' },
  modelValue: { type: String, default: '' },
  uri: { type: String }
})

const uriParsed = props.uri ? Uri.parse(props.uri) : undefined

const emit = defineEmits(['update:modelValue'])

const dialog = ref<boolean>(false)
const success = ref<boolean>(false)
const content = ref<string>(props.modelValue)

const update = (value: string) => {
  content.value = value
}

const save = () => {
  emit('update:modelValue', content.value)
  dialog.value = false
  success.value = true
}

let monaco: editor.IStandaloneCodeEditor | undefined = undefined

const initRefresh = (e: editor.IStandaloneCodeEditor) => {
  content.value = props.modelValue
  e.setValue(props.modelValue)

  e.addCommand(KeyMod.CtrlCmd | KeyCode.KeyS, () => {
    save()
  })

  e.addCommand(KeyMod.CtrlCmd | KeyCode.Escape, () => {
    dialog.value = false
  })

  monaco = e
}

watch(dialog, (open) => {
  if (!open) return

  content.value = props.modelValue
  monaco?.setValue(content.value)
  setTimeout(() => {
    monaco?.focus()
  }, 500)
})

const options: editor.IStandaloneEditorConstructionOptions = {
  colorDecorators: true,
  lineHeight: 24,
  tabSize: 2,
  minimap: { enabled: false }
}
</script>
