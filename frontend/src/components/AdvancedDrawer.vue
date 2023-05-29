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
          <AdvancedConfigDialog id="context" title="Context" :info="options.panels.contextInfo" v-model="inputs.context" uri="context.yaml" />
        </v-list-item>
        <v-list-item>
          <KyvernoConfig />
        </v-list-item>
        <v-list-item>
          <AdvancedConfigDialog id="crd" title="Custom Resource Definitions" :info="options.panels.crdInfo" v-model="inputs.customResourceDefinitions" />
        </v-list-item>
        <v-list-item>
          <AdvancedConfigDialog
            id="exceptions"
            title="Policy Exceptions"
            :info="options.panels.exceptionsInfo"
            v-model="inputs.exceptions"
            uri="policyexception.yaml" />
        </v-list-item>
        <v-list-item>
          <AdvancedConfigDialog id="clusterResources" title="Cluster Resources" :info="options.panels.clusterResourcesInfo" v-model="inputs.clusterResources">
            <template #append="{ update, content }">
              <ClusterResourceButton @update:model-value="update" :model-value="content" />
            </template>
          </AdvancedConfigDialog>
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

const props = defineProps({
  modelValue: { type: Boolean, default: false }
})

const emit = defineEmits(['update:modelValue'])
</script>
