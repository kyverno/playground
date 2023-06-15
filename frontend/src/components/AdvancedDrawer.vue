<template>
  <v-navigation-drawer
    :model-value="props.modelValue"
    temporary
    width="400"
    location="right"
    @update:modelValue="(event: boolean) => emit('update:modelValue', event)">
    <v-card flat>
      <v-list>
        <v-list-item>
          <AdvancedConfigDialog
            id="context"
            title="Context"
            :info="options.panels.contextInfo"
            v-model="inputs.context"
            uri="context.yaml"
            :template="ContextTemplate" />
        </v-list-item>
        <v-list-item>
          <KyvernoConfig />
        </v-list-item>
        <v-list-item>
          <AdvancedConfigDialog id="crd" title="Custom Resource Definitions" :info="options.panels.crdInfo" v-model="inputs.customResourceDefinitions">
            <template #actions="{ update, content }">
              <ClusterCRDButton @update:model-value="update" :model-value="content" label="From Cluster" v-if="config.cluster" />
            </template>
          </AdvancedConfigDialog>
        </v-list-item>
        <v-list-item>
          <AdvancedConfigDialog
            id="exceptions"
            title="Policy Exceptions"
            :info="options.panels.exceptionsInfo"
            v-model="inputs.exceptions"
            :template="PolicyExceptionTemplate"
            uri="policyexception.yaml">
            <template #actions="{ update, content }">
              <ClusterExceptionButton @update:model-value="update" :model-value="content" v-if="config.cluster" />
            </template>
          </AdvancedConfigDialog>
        </v-list-item>
        <v-list-item>
          <AdvancedConfigDialog id="clusterResources" title="Cluster Resources" :info="options.panels.clusterResourcesInfo" v-model="inputs.clusterResources">
            <template #actions="{ update, content }">
              <ClusterResourceButton @update:model-value="update" :model-value="content" label="From Cluster" v-if="config.cluster" />
            </template>
          </AdvancedConfigDialog>
        </v-list-item>
        <v-list-item>
          <AdvancedConfigDialog id="imageDate" title="Image Data" :info="options.panels.imageDataInfo" v-model="inputs.imageData" />
        </v-list-item>
      </v-list>
    </v-card>

    <template v-slot:append>
      <div class="pa-2">
        <v-btn flat color="primary" block @click="() => emit('update:modelValue', false)">Close</v-btn>
      </div>
    </template>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { options } from '@/config'
import AdvancedConfigDialog from './Advanced/AdvancedConfigDialog.vue'
import KyvernoConfig from './Advanced/KyvernoConfig.vue'
import { inputs } from '@/store'
import ClusterResourceButton from './Panel/ClusterResourceButton.vue'
import ClusterExceptionButton from './Panel/ClusterExceptionButton.vue'
import ClusterCRDButton from './Panel/ClusterCRDButton.vue'
import { ContextTemplate, PolicyExceptionTemplate } from '@/assets/templates'
import { config } from '@/composables/api'

const props = defineProps({
  modelValue: { type: Boolean, default: false }
})

const emit = defineEmits(['update:modelValue'])
</script>
