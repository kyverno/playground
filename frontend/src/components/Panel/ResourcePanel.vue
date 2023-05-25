<template>
  <v-card class="mt-3">
    <EditorToolbar
      title="Resources"
      id="resource-panel"
      :info="options.panels.resourceInfo"
      v-model="inputs.resource"
      :restore-value="state.resource.value"
    >
      <template #prepend-actions>
        <TemplateButton @select="onSelect" />
        <v-btn icon="mdi-format-page-split" @click="() => inputs.diffResources = !inputs.diffResources" />
      </template>
      <template #append-actions v-if="config.cluster">
        <ClusterResourceButton />
      </template>
    </EditorToolbar>
    <template v-if="inputs.diffResources">
      <v-toolbar density="compact" height="40" flat border color="grey-darken-3">
        <span class="d-inline-block pl-11 text-subtitle-2" style="width: 50%">Old Resource</span>
        <span class="d-inline-block pl-8 text-subtitle-2" style="width: 50%">New Resource</span>
      </v-toolbar>
      <ResourceDiffEditor v-model="inputs.resource" :height="height" />
    </template>
    <ResourceEditor v-model="inputs.resource" :height="height" v-else />
  </v-card>
</template>

<script setup lang="ts">
import { options } from '@/config';
import { useState } from '@/composables';
import TemplateButton from './TemplateButton.vue';
import ResourceEditor from './ResourceEditor.vue';
import ResourceDiffEditor from './ResourceDiffEditor.vue';
import EditorToolbar from './EditorToolbar.vue';
import ClusterResourceButton from './ClusterResourceButton.vue';
import { config } from '@/composables/api';
import { inputs } from '@/store';

const state = useState()

const onSelect = (template: string) => {
  inputs.resource = template

  if (inputs.diffResources) {
    inputs.oldResource = template
  }
}

defineProps({
    modelValue: { type: String, default: '' },
    height: { type: Number, default: 441 }
})

</script>