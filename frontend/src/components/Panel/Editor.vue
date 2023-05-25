<template>
    <MonacoEditor @editorDidMount="monacoSetup" :options="options" :modelValue="modelValue" @update:modelValue="(c: string) => emit('update:modelValue', c)" />
</template>

<script setup lang="ts">
import { PropType, reactive, watch } from 'vue';
import { KeyMod, KeyCode, editor } from "monaco-editor";
import MonacoEditor from './MonacoEditor.vue';

const props = defineProps({
    options: { type: Object as PropType<editor.IStandaloneEditorConstructionOptions>, default: () => ({}) },
    modelValue: { type: String, default: '' },
})

const emit = defineEmits(["editorDidMount", "update:modelValue"])
const options = reactive({ ...props.options })

const monacoSetup = (edit: editor.IStandaloneCodeEditor) => {
    edit.addCommand(KeyMod.Alt | KeyCode.KeyZ, () => {
        options.wordWrap = options.wordWrap !== 'on' ? 'on' : 'off'
    })

    watch(() => props.modelValue, (current, old) => {
        const currentLines = current.split(/\r\n|\r|\n/).length
        const oldLines = old.split(/\r\n|\r|\n/).length

        if (currentLines + 10 < oldLines) {
            edit.setScrollPosition({ scrollTop: 0 });
        }
    })

    emit('editorDidMount', edit)
}

</script>