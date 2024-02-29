import { useLocalStorage } from '@vueuse/core'
import { Ref, ref, watch } from 'vue'

export type Inputs = {
  name?: string
  diffResources?: boolean
  policy?: string
  oldResource?: string
  resource?: string
  context?: string
  config?: string
  exceptions?: string
  vapBindings?: string
  clusterResources?: string
  customResourceDefinitions?: string
  imageData?: string
}

const persisted = useLocalStorage<string>('persist:list', '')

export const getPersisted = (): Ref<string[]> => {
  const list = ref<string[]>([])

  watch(
    persisted,
    (content: string) => {
      list.value = (content || '').split(';;').filter((l) => !!l)
    },
    { immediate: true }
  )

  return list
}

export const createInput = (name: string, defaults?: Inputs) => {
  name = name.replaceAll(';;', ';').trim()
  const policy = useLocalStorage<string | undefined>(`persist:policy:${name}`, defaults?.policy || undefined)
  const resource = useLocalStorage<string | undefined>(`persist:resource:${name}`, defaults?.resource || undefined)
  const oldResource = useLocalStorage<string | undefined>(`persist:resource:old:${name}`, defaults?.oldResource || undefined)
  const context = useLocalStorage<string | undefined>(`persist:context:${name}`, defaults?.context || undefined)
  const config = useLocalStorage<string | undefined>(`persist:config:${name}`, defaults?.config || undefined)
  const exceptions = useLocalStorage<string | undefined>(`persist:exceptions:${name}`, defaults?.exceptions || undefined)
  const vapBindings = useLocalStorage<string | undefined>(`persist:vapBindings:${name}`, defaults?.vapBindings || undefined)
  const clusterResources = useLocalStorage<string | undefined>(`persist:clusterResources:${name}`, defaults?.clusterResources || undefined)
  const customResourceDefinitions = useLocalStorage<string | undefined>(`persist:crds:${name}`, defaults?.customResourceDefinitions || undefined)
  const imageData = useLocalStorage<string | undefined>(`persist:imageData:${name}`, defaults?.imageData || undefined)

  persisted.value = [...new Set([...getPersisted().value, name])].join(';;')

  return {
    policy,
    resource,
    oldResource,
    context,
    config,
    exceptions,
    vapBindings,
    clusterResources,
    customResourceDefinitions,
    imageData,
    name
  }
}

export const updateInput = (name: string, values: Inputs) => {
  const input = createInput(name)

  if (input.policy.value !== values.policy) {
    input.policy.value = values.policy
  }
  if (input.resource.value !== values.resource) {
    input.resource.value = values.resource
  }
  if (input.oldResource.value !== values.oldResource) {
    input.oldResource.value = values.oldResource
  }
  if (input.context.value !== values.context) {
    input.context.value = values.context
  }
  if (input.exceptions.value !== values.exceptions) {
    input.exceptions.value = values.exceptions
  }
  if (input.vapBindings.value !== values.vapBindings) {
    input.vapBindings.value = values.vapBindings
  }
  if (input.config.value !== values.config) {
    input.config.value = values.config
  }
  if (input.clusterResources.value !== values.clusterResources) {
    input.clusterResources.value = values.clusterResources
  }
  if (input.customResourceDefinitions.value !== values.customResourceDefinitions) {
    input.customResourceDefinitions.value = values.customResourceDefinitions
  }
  if (input.imageData.value !== values.imageData) {
    input.imageData.value = values.imageData
  }

  return input
}

export const removeInput = (name: string) => {
  const input = createInput(name)

  input.policy.value = null
  input.resource.value = null
  input.oldResource.value = null
  input.context.value = null
  input.config.value = null
  input.exceptions.value = null
  input.vapBindings.value = null
  input.clusterResources.value = null
  input.customResourceDefinitions.value = null
  input.imageData.value = null

  name = name.replaceAll(';;', ';').trim()
  const list = getPersisted()

  persisted.value = list.value.filter((l) => l !== name).join(';;')
}
