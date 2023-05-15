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
import { options } from '@/config';
import ContextEditor from './ContextEditor.vue';
import EditorToolbar from './EditorToolbar.vue';
import { useState } from '@/composables';
import { ref } from 'vue';
import { computed } from 'vue';

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