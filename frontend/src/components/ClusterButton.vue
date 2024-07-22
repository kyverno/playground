<template>
  <v-dialog v-model="dialog" width="600px" :theme="layoutTheme">
    <template v-slot:activator="{ props }">
      <v-tooltip location="bottom" content-class="no-opacity-tooltip" text="Import from Cluster" theme="dark" v-if="!label">
        <template v-slot:activator="{ props: tooltip }">
          <v-btn v-bind="{ ...tooltip, ...props }" icon="mdi-kubernetes" />
        </template>
      </v-tooltip>
      <v-btn v-bind="props" prepend-icon="mdi-kubernetes" v-else>{{ label }}</v-btn>
    </template>

    <v-card :theme="layoutTheme" title="Load from connected Cluster">
      <v-divider class="my-2" />
      <v-window v-model="window">
        <v-window-item>
          <v-container>
            <simple-row v-if="error || errorResources">
              <v-alert color="error" variant="outlined">{{ error || errorResources }}</v-alert>
            </simple-row>
            <slot name="resource-api" :update="updateResource" :resource="resourceAPI"></slot>
            <simple-row v-if="!resourceAPI.clusterScoped">
              <namespace-select v-model="namespace" />
            </simple-row>
            <simple-row>
              <v-text-field v-model="name" label="Name" variant="outlined" hide-details density="comfortable" />
            </simple-row>
            <simple-row>
              <mode-select v-model="mode" />
            </simple-row>
          </v-container>
          <v-card-actions>
            <v-btn @click="dialog = false">Close</v-btn>
            <v-spacer />
            <v-btn @click="search" :loading="loadingResources" :color="errorResources ? 'error' : undefined">Search</v-btn>
            <v-btn
              :loading="loading"
              :disabled="!name || !resourceAPI"
              @click="() => load([{ name: name || '', namespace }])"
              :color="error ? 'error' : undefined">
              Load Resource
            </v-btn>
          </v-card-actions>
        </v-window-item>
        <v-window-item>
          <v-container>
            <simple-row v-if="error">
              <v-alert color="error" variant="outlined">{{ error }}</v-alert>
            </simple-row>
            <simple-row>
              <cluster-search-list v-model="selections" :foundings="foundings" />
            </simple-row>
            <simple-row>
              <mode-select v-model="mode" />
            </simple-row>
          </v-container>
          <v-card-actions>
            <v-btn @click="dialog = false">Close</v-btn>
            <v-btn @click="window = 0">Back</v-btn>
            <v-spacer />
            <v-btn :loading="loading" :disabled="!selections.length || !resourceAPI" @click="() => load(selections)" :color="error ? 'error' : undefined">
              Load Resources
            </v-btn>
          </v-card-actions>
        </v-window-item>
      </v-window>
    </v-card>
  </v-dialog>
</template>
<script setup lang="ts">
import { ref, watch, type PropType } from 'vue'
import { layoutTheme } from '@/config'
import { useAPI, resourcesToYAML, type ResourceKind } from '@/composables/api'
import { mergeResources } from '@/utils'
import ClusterSearchList from '@/components/Panel/ClusterSearchList.vue'
import ModeSelect from '@/components/ModeSelect.vue'
import NamespaceSelect from '@/components/NamespaceSelect.vue'
import SimpleRow from '@/components/SimpleRow.vue'

type Resource = { namespace?: string; name: string }
type ResourceAPI = ResourceKind & { title: string }

const props = defineProps({
  modelValue: { type: String },
  label: { type: String },
  defaultResource: { type: Object as PropType<ResourceAPI>, required: true }
})

const emit = defineEmits(['update:modelValue'])

const window = ref<number>(0)
const dialog = ref<boolean>(false)

const resourceAPI = ref<ResourceAPI>(props.defaultResource)
const updateResource = (res: ResourceAPI) => {
  resourceAPI.value = res
}

const namespace = ref<string>()
const name = ref<string>()
const mode = ref<string>('replace')

watch(resourceAPI, (api) => {
  if (api.clusterScoped) {
    namespace.value = undefined
  }
})

const { loading: loadingResources, error: errorResources, resources: loadResources, data: foundings } = useAPI<Resource[]>()

const search = () => {
  const { apiVersion, kind } = resourceAPI.value

  loadingResources.value = true
  loadResources({ apiVersion, kind, namespace: namespace.value || '' })
    .then((resources) => {
      foundings.value = resources

      window.value = 1
    })
    .catch((err) => {
      errorResources.value = err
    })
    .finally(() => {
      loadingResources.value = false
    })
}

const selections = ref<Resource[]>([])

const { loading, error, resource: loadResource } = useAPI<object[]>()

const load = (res: Resource[]) => {
  const { apiVersion, kind } = resourceAPI.value

  const promises = res.map(({ namespace, name }) => loadResource({ apiVersion, kind, namespace: namespace || '', name }))

  loading.value = true
  Promise.all(promises)
    .then((response) => {
      const results = resourcesToYAML(response)

      if (mode.value === 'append') {
        emit('update:modelValue', mergeResources(props.modelValue || '', results))
      } else {
        emit('update:modelValue', results)
      }
      dialog.value = false
    })
    .catch((err) => {
      error.value = err
    })
    .finally(() => {
      loading.value = false
    })
}

watch(dialog, (value: boolean) => {
  if (value) return

  setTimeout(() => {
    window.value = 0
    namespace.value = undefined
    name.value = undefined
    selections.value = []

    loading.value = false
    loadingResources.value = false

    error.value = undefined
    errorResources.value = undefined
  }, 300)
})
</script>
