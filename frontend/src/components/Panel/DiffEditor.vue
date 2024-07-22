<template>
  <div ref="root" class="monaco-editor-vue3" :style="style" :id="id"></div>
</template>

<script setup lang="ts">
import { ref, toRefs, computed, onMounted, reactive } from 'vue'
import * as monaco from 'monaco-editor'
import { type PropType } from 'vue'
import { watch } from 'vue'
import { onBeforeUnmount } from 'vue'
import { useEditorFix } from '@/functions/editor'

const props = defineProps({
  uri: { type: Object as PropType<monaco.Uri> },
  width: { type: [String, Number], default: '100%' },
  height: { type: [String, Number], default: '100%' },
  modelValue: { type: String, default: '' },
  original: { type: String, default: '' },
  language: { type: String, default: 'javascript' },
  theme: { type: String, default: 'vs' },
  id: { type: String, required: true },
  options: {
    type: Object as PropType<monaco.editor.IStandaloneDiffEditorConstructionOptions>,
    default: () => ({})
  }
})

const emit = defineEmits([
  'editorWillMount',
  'editorDidMount',
  'update:modelValue',
  'update:original',
  'switchWordWrap'
])

const { width, height } = toRefs(props)

const style = computed((): { [key: string]: string } => {
  const fixedWidth = typeof width.value === 'string' ? width.value : `${width.value}px`
  const fixedHeight = typeof height.value === 'string' ? height.value : `${height.value}px`
  return { width: fixedWidth, height: fixedHeight, 'text-align': 'left' }
})

const root = ref<HTMLElement>()
const options = reactive({ ...props.options })
let editor: monaco.editor.IStandaloneDiffEditor | undefined = undefined

onMounted(() => {
  emit('editorWillMount', monaco)

  let model: monaco.editor.ITextModel | null = null
  if (props.uri) {
    model = monaco.editor.getModel(props.uri)
  }
  if (!model) {
    model = monaco.editor.createModel(props.modelValue, props.language, props.uri)
  }

  if (!root.value) return

  editor = monaco.editor.createDiffEditor(root.value, {
    theme: props.theme,
    automaticLayout: true,
    ...props.options
  })

  editor.setModel({
    original: monaco.editor.createModel(props.original, props.language),
    modified: monaco.editor.createModel(props.modelValue, props.language)
  })

  const modifiedEditor = editor.getModifiedEditor()
  const originalEditor = editor.getOriginalEditor()

  const disposeModified = useEditorFix(modifiedEditor, `modified-${props.id}`)
  const disposeOriginal = useEditorFix(originalEditor, `original-${props.id}`)

  modifiedEditor.onDidChangeModelContent(() => {
    const value = modifiedEditor.getValue()
    if (props.modelValue !== value) {
      emit('update:modelValue', value)
    }
  })

  originalEditor.onDidChangeModelContent(() => {
    const value = originalEditor.getValue()
    if (props.original !== value) {
      emit('update:original', value)
    }
  })

  watch(
    () => props.modelValue,
    (current, old) => {
      const currentLines = current.split(/\r\n|\r|\n/).length
      const oldLines = old.split(/\r\n|\r|\n/).length

      if (currentLines + 10 < oldLines) {
        modifiedEditor?.setScrollPosition({ scrollTop: 0 })
      }

      if (current !== modifiedEditor.getValue()) {
        modifiedEditor.setValue(current)
      }
    }
  )

  watch(
    () => props.original,
    (current, old) => {
      const currentLines = current.split(/\r\n|\r|\n/).length
      const oldLines = old.split(/\r\n|\r|\n/).length

      if (currentLines + 10 < oldLines) {
        originalEditor?.setScrollPosition({ scrollTop: 0 })
      }

      if (current !== originalEditor.getValue()) {
        originalEditor.setValue(current)
      }
    }
  )

  modifiedEditor.addCommand(monaco.KeyMod.Alt | monaco.KeyCode.KeyZ, () => {
    options.wordWrap = options.wordWrap !== 'on' ? 'on' : 'off'
  })

  originalEditor.addCommand(monaco.KeyMod.Alt | monaco.KeyCode.KeyZ, () => {
    options.wordWrap = options.wordWrap !== 'on' ? 'on' : 'off'
  })

  onBeforeUnmount(() => {
    disposeModified()
    disposeOriginal()
    editor?.dispose()
  })

  emit('editorDidMount', editor)
})

watch(
  options,
  (o) => {
    editor?.updateOptions({ ...o })
  },
  { deep: true }
)

watch(
  () => props.options,
  (o: monaco.editor.IEditorOptions) => {
    for (const config in o) {
      // @ts-ignore
      options[config] = o[config]
    }
  },
  { deep: true }
)

watch(
  () => props.theme,
  (theme) => {
    monaco.editor.setTheme(theme)
  }
)
</script>
