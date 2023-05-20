<template>
    <v-autocomplete 
        variant="outlined" 
        hide-details 
        :items="list" 
        label="Policy Type"
        :model-value="modelValue"
        @update:model-value="e => emit('update:modelValue', e)"
        return-object
        density="comfortable"
    />
</template>

<script setup lang="ts">
import { PropType } from 'vue';
import { watch } from 'vue';

const list = [
    { title: 'kyverno.io/v1 ClusterPolicy', apiVersion: 'kyverno.io/v1', kind: 'ClusterPolicy', clusterScoped: true },
    { title: 'kyverno.io/v1 Policy', apiVersion: 'kyverno.io/v1', kind: 'Policy', clusterScoped: false },
    { title: 'kyverno.io/v2beta1 ClusterPolicy', apiVersion: 'kyverno.io/v2beta1', kind: 'ClusterPolicy', clusterScoped: true },
    { title: 'kyverno.io/v2beta1 Policy', apiVersion: 'kyverno.io/v2beta1', kind: 'Policy', clusterScoped: false },
]

const props = defineProps({
    modelValue: { type: Object as PropType<{ title: string; apiVersion: string; kind: string; clusterScoped: boolean }> }
})

const emit = defineEmits(['update:modelValue'])

watch(() => props.modelValue, (v) => { if (!v) emit('update:modelValue', list[0]) })
</script>