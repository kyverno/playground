<template>
  <AdvancedConfigDialog id="config" title="Kyverno Config" :info="options.panels.configInfo" v-model="inputs.config">
    <template #actions="{ update }">
      <v-btn @click="() => update('')">Clear Config</v-btn>
      <v-btn @click="() => update(ConfigTemplate)">Load Default Config</v-btn>
      <UploadButton label="Upload ConfigMap" @click="update" :tooltip="false" />
      <ClusterConfigButton v-if="config.cluster" @update="update" />
    </template>
  </AdvancedConfigDialog>
</template>

<script setup lang="ts">
import { watch } from 'vue'
import { inputs } from '@/store'
import UploadButton from '@/components/Panel/UploadButton.vue'
import ClusterConfigButton from '@/components/ClusterConfigButton.vue'
import { ConfigTemplate } from '@/assets/templates'
import { config } from '@/composables/api'
import { loadedConfig } from '@/composables'
import AdvancedConfigDialog from './AdvancedConfigDialog.vue'
import { options } from '@/config'

watch(
  () => inputs.config,
  (config?: string) => {
    loadedConfig.value = config || ''
  }
)
</script>
