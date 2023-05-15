<template>
  <v-card class="mt-3">
    <EditorToolbar
      title="Resources"
      id="resource-panel"
      :info="options.panels.resourceInfo"
      v-model="resource"
      :restore-value="state.resource.value"
    >
      <template #prepend-actions>
        <TemplateButton @select="(template: string) => resource = template" />
      </template>
    </EditorToolbar>
    <ResourceEditor v-model="resource" :height="height" />
  </v-card>
</template>

<script setup lang="ts">
import { options } from '@/config';
import ResourceEditor from './ResourceEditor.vue';
import EditorToolbar from './EditorToolbar.vue';
import { useState } from '@/composables';
import { computed } from 'vue';
import { TemplateButton } from './buttons';

const state = useState()

const props = defineProps({
    modelValue: { type: String, default: '' },
    height: { type: Number, default: 441 }
})

const emit = defineEmits(["update:modelValue", "collapse"])

const resource = computed({
    get() {
        return props.modelValue
    },
    set(value: string) {
        emit('update:modelValue', value)
    }
})

</script>