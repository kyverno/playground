<template>
<v-card style="height: 800px">
  <EditorToolbar
    id="policy-panel"
    title="Policies"
    v-model="policy"
    :restore-value="state.policy.value"
    :info="options.panels.policyInfo"
  >
    <template #append-actions v-if="config.cluster">
      <ClusterPolicyButton />
    </template>
  </EditorToolbar>
  <PolicyEditor v-model="policy" />
</v-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { options } from '@/config';
import { useState } from '@/composables';
import PolicyEditor from './PolicyEditor.vue';
import EditorToolbar from './EditorToolbar.vue';
import { config } from '@/composables/api';
import ClusterPolicyButton from './ClusterPolicyButton.vue';

const state = useState()

const props = defineProps({
    modelValue: { type: String, default: '' }
})

const emit = defineEmits(["update:modelValue", "collapse"])

const policy = computed({
    get() {
        return props.modelValue
    },
    set(value: string) {
        emit('update:modelValue', value)
    }
})
</script>