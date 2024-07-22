<template>
  <v-card>
    <EditorToolbar
      id="context-panel"
      title="Context"
      :info="options.panels.contextInfo"
      v-model="inputs.context"
      :restore-value="state.context.value"
    >
      <template #append-actions>
        <v-btn :icon="collapse ? 'mdi-chevron-up' : 'mdi-chevron-down'" @click="toggle" />
      </template>
    </EditorToolbar>
    <ContextEditor v-model="inputs.context" v-show="!collapse" style="height: 250px" />
  </v-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { options } from '@/config'
import { useState } from '@/composables'
import ContextEditor from './ContextEditor.vue'
import EditorToolbar from './EditorToolbar.vue'
import { inputs } from '@/store'

const state = useState()
const collapse = ref(false)

const emit = defineEmits(['collapse'])

const toggle = () => {
  collapse.value = !collapse.value
  emit('collapse', collapse.value)
}
</script>
