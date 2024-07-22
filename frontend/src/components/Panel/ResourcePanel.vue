<template>
  <v-card :height="height">
    <EditorToolbar
      title="Resources"
      id="resource-panel"
      :info="options.panels.resourceInfo"
      v-model="inputs.resource"
      :restore-value="state.resource.value"
    >
      <template #prepend-actions>
        <TemplateButton @select="onSelect" />
        <v-tooltip
          location="bottom"
          content-class="no-opacity-tooltip"
          text="Switch between normal and difference mode to define an old and new resource"
        >
          <template v-slot:activator="{ props }">
            <v-btn icon="mdi-format-page-split" v-bind="props" @click="switchEditor" />
          </template>
        </v-tooltip>
      </template>
      <template #append-actions v-if="config.cluster">
        <ClusterResourceButton @update:model-value="onSelect" :model-value="inputs.resource" />
      </template>
    </EditorToolbar>
    <template v-if="inputs.diffResources">
      <v-toolbar
        density="compact"
        height="40"
        flat
        border
        color="grey-darken-3"
        v-show="inputs.diffResources"
      >
        <span class="d-inline-block pl-11 text-subtitle-2" style="width: 50%">Old Resource</span>
        <span class="d-inline-block pl-8 text-subtitle-2" style="width: 50%">New Resource</span>
      </v-toolbar>
      <ResourceDiffEditor v-model="inputs.resource" />
    </template>
    <ResourceEditor v-model="inputs.resource" v-show="!inputs.diffResources" />
  </v-card>
</template>

<script setup lang="ts">
import { options } from '@/config'
import { useState } from '@/composables'
import TemplateButton from './TemplateButton.vue'
import ResourceEditor from './ResourceEditor.vue'
import ResourceDiffEditor from './ResourceDiffEditor.vue'
import EditorToolbar from './EditorToolbar.vue'
import ClusterResourceButton from './ClusterResourceButton.vue'
import { config } from '@/composables/api'
import { inputs } from '@/store'

const state = useState()

const onSelect = (template: string) => {
  inputs.resource = template

  if (inputs.diffResources) {
    inputs.oldResource = template
  }
}

const switchEditor = () => {
  inputs.diffResources = !inputs.diffResources
}

defineProps({
  modelValue: { type: String, default: '' },
  height: { type: Number, default: 441 }
})
</script>
