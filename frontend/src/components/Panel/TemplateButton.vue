<template>
  <v-menu location="bottom" :close-on-content-click="false" v-model="menu">
    <template #activator="{ props }">
      <v-tooltip
        location="top"
        content-class="no-opacity-tooltip"
        text="Load Resource from predefined Templates"
        theme="dark"
      >
        <template #activator="{ props: tooltip }">
          <v-btn v-bind="{ ...tooltip, ...props }" icon="mdi-note-text"></v-btn>
        </template>
      </v-tooltip>
    </template>
    <v-list>
      <v-list-item>
        <v-text-field
          variant="outlined"
          density="compact"
          hide-details
          placeholder="Search"
          class="pb-2"
          v-model="search"
        />
      </v-list-item>
      <v-divider />
      <v-virtual-scroll :items="templates" :height="350" width="210">
        <template v-slot:default="{ item }">
          <v-list-item @click="() => loadTemplate(item)" :title="item" />
        </template>
      </v-virtual-scroll>
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const emit = defineEmits(['select'])
const menu = ref(false)
const search = ref('')

const list = [
  'Pod',
  'Deployment',
  'StatefulSet',
  'DaemonSet',
  'ReplicaSet',
  'Job',
  'CronJob',
  'Service',
  'Ingress',
  'NetworkPolicy',
  'Namespace',
  'ConfigMap',
  'Secret',
  'Role',
  'RoleBinding',
  'ClusterRole',
  'ClusterRoleBinding',
  'PersistentVolume',
  'PersistentVolumeClaim',
  'ServiceAccount',
  'PodDisruptionBudget',
  'PriorityClass',
  'ResourceQuota'
]

const templates = computed(() => {
  if (!search.value) return list

  return list.filter((i) => i.toLowerCase().search(search.value.toLowerCase()) !== -1)
})

const loadTemplate = (template: string) => {
  return fetch(`templates/${template.toLocaleLowerCase()}.yaml`)
    .then((resp) => {
      if (resp.status !== 200) return

      return resp.text().then((content) => emit('select', content))
    })
    .then(() => {
      menu.value = false
      search.value = ''
    })
    .catch(console.error)
}
</script>
