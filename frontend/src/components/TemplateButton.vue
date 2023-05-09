<template>
  <v-menu location="bottom">
    <template v-slot:activator="{ props }">
      <v-btn v-bind="props" prepend-icon="mdi-note-text">from temlate</v-btn>
    </template>
    <v-list>
      <v-list-item v-for="(item, key) in templates" :key="key" @click="emit('select', item.template)" :title="item.title" />
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
import { PodTemplate, DeploymentTemplate, ServiceTemplate, IngressTemplate, ConfigMapTemplate, SecretTemplate } from "@/assets/templates";

const emit = defineEmits(['select'])

const templates = [
    { title: 'Pod', template: PodTemplate },
    { title: 'Deployment', template: DeploymentTemplate },
    { title: 'Service', template: ServiceTemplate },
    { title: 'Ingress', template: IngressTemplate },
    { title: 'ConfigMap', template: ConfigMapTemplate },
    { title: 'Secret', template: SecretTemplate },
]
</script>
