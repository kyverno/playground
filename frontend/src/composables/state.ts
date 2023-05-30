import { ContextTemplate, PolicyTemplate, ResourceTemplate, ConfigTemplate } from '@/assets/templates'
import { useLocalStorage } from '@vueuse/core'

export const loadedPolicy = useLocalStorage<string>('loaded:policy', PolicyTemplate)
export const loadedContext = useLocalStorage<string>('loaded:context', ContextTemplate)
export const loadedResource = useLocalStorage<string>('loaded:resource', ResourceTemplate)
export const loadedOldResource = useLocalStorage<string>('loaded:resource:old', '')
export const loadedConfig = useLocalStorage<string>('loaded:config', ConfigTemplate)
export const loadedExceptions = useLocalStorage<string>('loaded:exceptions', '')
export const loadedClusterResources = useLocalStorage<string>('loaded:clusterResources', '')
export const loadedCustomResourceDefinitions = useLocalStorage<string>('loaded:crds', '')
export const loadedState = useLocalStorage<string>('loaded:state', '')

export type State = {
  name: string
  policy?: string
  resource?: string
  oldResource?: string
  context?: string
  config?: string
  exceptions?: string
  clusterResources?: string
  customResourceDefinitions?: string
}

const update = (values: State) => {
  if (values.policy) {
    loadedPolicy.value = values.policy
  }

  if (values.resource) {
    loadedResource.value = values.resource
  }

  if (values.oldResource) {
    loadedOldResource.value = values.oldResource
  }

  if (values.context) {
    loadedContext.value = values.context
  }

  if (values.config) {
    loadedConfig.value = values.config
  }

  if (values.exceptions) {
    loadedExceptions.value = values.exceptions
  }

  if (values.clusterResources) {
    loadedClusterResources.value = values.clusterResources
  }

  if (values.customResourceDefinitions) {
    loadedCustomResourceDefinitions.value = values.customResourceDefinitions
  }

  loadedState.value = values.name
}

export const useState = () => ({
  config: loadedConfig,
  policy: loadedPolicy,
  resource: loadedResource,
  oldResource: loadedOldResource,
  context: loadedContext,
  exceptions: loadedExceptions,
  clusterResources: loadedClusterResources,
  customResourceDefinitions: loadedCustomResourceDefinitions,
  name: loadedState,
  update
})
