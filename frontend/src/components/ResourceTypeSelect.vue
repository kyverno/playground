<template>
    <v-autocomplete 
        variant="outlined" 
        hide-details 
        :items="list" 
        label="Resource Type"
        v-model="resource"
        return-object
        density="comfortable"
    />
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { PropType } from 'vue';

const list = [
  'v1 Pod',
  'v1 Service',
  'v1 Namespace',
  'v1 ConfigMap',
  'v1 Secret',
  'v1 PersistedVolume',
  'v1 PersistedVolumeClaim',
  'v1 ServiceAccount',
  'v1 ResourceQuota',
  'apps/v1 Deployment',
  'apps/v1 StatefulSet',
  'apps/v1 DaemonSet',
  'apps/v1 ReplicaSet',
  'batch/v1 Job',
  'batch/v1 CronJob',
  'networking.k8s.io/v1 Ingress',
  'networking.k8s.io/v1 NetworkPolicy',
  'rbac.authorization.k8s.io/v1 Role',
  'rbac.authorization.k8s.io/v1 RoleBinding',
  'rbac.authorization.k8s.io/v1 ClusterRole',
  'rbac.authorization.k8s.io/v1 ClusterRoleBinding',
  'policy/v1 PodDisruptionBudget',
  'scheduling.k8s.io/v1 PriorityClass',
]

const props = defineProps({
    modelValue: { type: Object as PropType<{ apiVersion: string; kind: string }>, required: true }
})

const emit = defineEmits(['update:modelValue'])

const resource = computed({
    get() {
        return `${props.modelValue.apiVersion} ${props.modelValue.kind}` 
    },
    set(val: string) {
        const parts = val.split(' ')
        
        emit('update:modelValue', { apiVersion: parts[0] || '', kind: parts[1] || '' })
    }
})
</script>