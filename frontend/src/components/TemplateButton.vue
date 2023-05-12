<template>
  <v-menu location="bottom">
    <template v-slot:activator="{ props }">
      <v-btn v-bind="props" prepend-icon="mdi-note-text">from template</v-btn>
    </template>
    <v-list>
      <v-list-item v-for="(item, key) in templates" :key="key" @click="() => loadTemplate(item)" :title="item" />
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">

const emit = defineEmits(['select'])

const templates = [
  'Pod',
  'Deployment',
  'Service',
  'Ingress',
  'ConfigMap',
  'Secret'
]

const loadTemplate = (template: string) => {
  return fetch(`/templates/${template.toLocaleLowerCase()}.yaml`).then((resp) => {
    if (resp.status !== 200) return

    return resp.text().then(content => emit('select', content))
  }).catch(console.error)
}
</script>
