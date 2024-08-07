import {
  ContextTemplate,
  PolicyTemplate,
  ResourceTemplate,
  ConfigTemplate
} from '@/assets/templates'
import { useLocalStorage } from '@vueuse/core'

export const loadedPolicy = useLocalStorage<string>('loaded:policy', PolicyTemplate)
export const loadedContext = useLocalStorage<string>('loaded:context', ContextTemplate)
export const loadedResource = useLocalStorage<string>('loaded:resource', ResourceTemplate)
export const loadedOldResource = useLocalStorage<string>('loaded:resource:old', '')
export const loadedConfig = useLocalStorage<string>('loaded:config', ConfigTemplate)
export const loadedExceptions = useLocalStorage<string>('loaded:exceptions', '')
export const loadedVapBindings = useLocalStorage<string>('loaded:vapBindings', '')
export const loadedClusterResources = useLocalStorage<string>('loaded:clusterResources', '')
export const loadedCustomResourceDefinitions = useLocalStorage<string>('loaded:crds', '')
export const loadedImageData = useLocalStorage<string>('loaded:imageData', '')
export const loadedState = useLocalStorage<string>('loaded:state', '')

export type State = {
  name: string
  policy?: string
  resource?: string
  oldResource?: string
  context?: string
  config?: string
  exceptions?: string
  vapBindings?: string
  clusterResources?: string
  customResourceDefinitions?: string
  imageData?: string
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

  if (values.vapBindings) {
    loadedVapBindings.value = values.vapBindings
  }

  if (values.clusterResources) {
    loadedClusterResources.value = values.clusterResources
  }

  if (values.customResourceDefinitions) {
    loadedCustomResourceDefinitions.value = values.customResourceDefinitions
  }

  if (values.imageData) {
    loadedImageData.value = values.imageData
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
  vapBindings: loadedVapBindings,
  clusterResources: loadedClusterResources,
  customResourceDefinitions: loadedCustomResourceDefinitions,
  imageData: loadedImageData,
  name: loadedState,
  update
})
