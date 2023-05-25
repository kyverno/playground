<template>
  <div ref="root" class="monaco-editor-vue3" :style="style"></div>
</template>

<script setup lang="ts">
import { ref, toRefs, computed, onMounted } from 'vue';
import * as monaco from "monaco-editor";
import { PropType } from 'vue';
import { watch } from 'vue';

const props = defineProps({
    uri: { type: Object as PropType<monaco.Uri> },
    width: { type: [String, Number], default: "100%" },
    height: { type: [String, Number], default: "100%" },
    modelValue: { type: String, default: '' },
    language: { type: String, default: "javascript" },
    theme: { type: String, default: "vs" },
    options: { type: Object as PropType<monaco.editor.IStandaloneEditorConstructionOptions>, default: () => ({}) },
})

const emit = defineEmits(["editorWillMount", "editorDidMount", "update:modelValue", "switchWordWrap"])

const { width, height } = toRefs(props);

const style = computed((): { [key: string]: string } => {
  const fixedWidth = typeof width.value === 'string'  ? width.value : `${width.value}px`;
  const fixedHeight = typeof height.value === 'string' ? height.value : `${height.value}px`;
  return { width: fixedWidth, height: fixedHeight, "text-align": "left" };
});

const root = ref<HTMLElement>()
let editor: monaco.editor.IStandaloneCodeEditor | undefined = undefined;

onMounted(() => {
  emit("editorWillMount", monaco);

  let model: monaco.editor.ITextModel | null = null;
  if (props.uri) {
    model = monaco.editor.getModel(props.uri);
  }
  if (!model) {
    model = monaco.editor.createModel(props.modelValue, props.language, props.uri);
  }

  if (!root.value) return;

  editor = monaco.editor.create(root.value, {
    value: props.modelValue,
    language: props.language,
    theme: props.theme,
    automaticLayout: true,
    ...props.options,
    model,
  });

  editor.onDidChangeModelContent(() => {
    const value = editor?.getValue();
    if (props.modelValue !== value) {
      emit("update:modelValue", value);
    }
  });

  emit("editorDidMount", editor);
})

watch(() => props.options, (o) => {
  editor?.updateOptions({ ...o });
}, { deep: true })

watch(() => props.modelValue, (value) => {
  if (value === editor?.getValue()) return

  editor?.setValue(value)
})

watch(() => props.language, (language) => {
  if (!editor) return;
  const model = editor.getModel()  
  if (!model) return;

  monaco.editor.setModelLanguage(model, language);
})

watch(() => props.theme, (theme) => {
  monaco.editor.setTheme(theme);
})
</script>
