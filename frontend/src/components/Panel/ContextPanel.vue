<template>
    <v-card>
      <EditorToolbar
        id="context-panel"
        title="Context"
        :info="options.panels.contextInfo"
        v-model="context"
        :restore-value="state.context.value"
      >
        <template #append-actions>
          <v-btn :icon="collapse ? 'mdi-chevron-up' : 'mdi-chevron-down'" @click="toggle" />
        </template>
      </EditorToolbar>
      <ContextEditor v-model="context" v-show="!collapse" style="height: 250px" />
    </v-card>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { options } from '@/config';
import { useState } from '@/composables';
import ContextEditor from './ContextEditor.vue';
import EditorToolbar from './EditorToolbar.vue';

const state = useState()
const collapse = ref(false)

const props = defineProps({
    modelValue: { type: String, default: '' }
})

const emit = defineEmits(["update:modelValue", "collapse"])

const context = computed({
    get() {
        return props.modelValue
    },
    set(value: string) {
        emit('update:modelValue', value)
    }
})

const toggle = () => {
    collapse.value = !collapse.value
    emit('collapse', collapse.value)
}

</script>