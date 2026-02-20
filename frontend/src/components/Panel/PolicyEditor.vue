<template>
  <div style="position: relative">
    <MonacoEditor
      language="yaml"
      :theme="editorTheme"
      :modelValue="props.modelValue"
      @update:modelValue="(event: string) => emit('update:modelValue', event)"
      :options="options"
      @editorDidMount="monacoSetup"
      :uri="uri"
      id="policy"
      :height="752"
    />
    <v-card class="config" theme="dark" color="black" v-if="false">
      <v-card-text class="my-0 py-1">
        <v-switch
          v-model="autocompleteOnEnter"
          label="autocomplete"
          hide-details
          density="compact"
        />
      </v-card-text>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { editor, KeyCode, Uri } from 'monaco-editor'
import { editorTheme } from '@/config'
import { loadedPolicy } from '@/composables'
import MonacoEditor from './MonacoEditor.vue'

const props = defineProps({
  modelValue: { type: String, default: '' },
})

const emit = defineEmits(['update:modelValue'])

const uri = Uri.parse('policy.yaml')

const autocompleteOnEnter = ref(true)
const eventRegistered = ref(false)

const options: editor.IStandaloneEditorConstructionOptions = {
  wordWrap: 'off',
  colorDecorators: true,
  lineHeight: 24,
  tabSize: 2,
}

let timeout: any = undefined

const monacoSetup = (edit: editor.IStandaloneCodeEditor) => {
  watch(loadedPolicy, () => {
    edit.setScrollPosition({ scrollTop: 0 })
  })

  if (eventRegistered.value) return

  edit.onKeyUp((e) => {
    if (!autocompleteOnEnter.value) return

    const position = edit.getPosition()
    if (!position) return

    const text = edit.getModel()?.getLineContent(position.lineNumber).trim()

    if (e.keyCode === KeyCode.Enter && !text) {
      if (typeof timeout === 'number') {
        clearTimeout(timeout)
      }

      timeout = setTimeout(() => {
        edit.trigger('', 'editor.action.triggerSuggest', '')
      }, 1500)
    } else if (e.keyCode === KeyCode.Escape && typeof timeout === 'number') {
      clearTimeout(timeout)
      timeout = undefined
    } else if (typeof timeout === 'number' && !text) {
      clearTimeout(timeout)

      timeout = setTimeout(() => {
        edit.trigger('', 'editor.action.triggerSuggest', '')
      }, 1500)
    } else if (typeof timeout === 'number') {
      clearTimeout(timeout)
      timeout = undefined
    }
  })

  eventRegistered.value = true
}
</script>

<style scoped>
.config {
  position: absolute;
  bottom: 60px;
  width: 80%;
  margin-left: 10%;
}
</style>
