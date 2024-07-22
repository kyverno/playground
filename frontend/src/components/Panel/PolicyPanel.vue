<template>
  <v-card style="height: 800px">
    <EditorToolbar
      id="policy-panel"
      title="Policies"
      v-model="inputs.policy"
      :restore-value="state.policy.value"
      :info="options.panels.policyInfo"
    >
      <template #append-actions v-if="config.cluster">
        <ClusterPolicyButton />
      </template>
    </EditorToolbar>
    <PolicyEditor v-model="inputs.policy" />
  </v-card>
</template>

<script setup lang="ts">
import { options } from '@/config'
import { useState } from '@/composables'
import PolicyEditor from './PolicyEditor.vue'
import EditorToolbar from './EditorToolbar.vue'
import { config } from '@/composables/api'
import ClusterPolicyButton from './ClusterPolicyButton.vue'
import { inputs } from '@/store'

const state = useState()
</script>
