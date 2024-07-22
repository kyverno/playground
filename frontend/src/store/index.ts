import {
  ConfigTemplate,
  ContextTemplate,
  PolicyTemplate,
  ResourceTemplate,
  CustomResourceDefinitionsTemplate,
  ImageDataTemplate
} from '@/assets/templates'
import { reactive } from 'vue'
import { useState, type Inputs } from '@/composables'
import { type ResourceKind } from '@/composables/api'

type State = {
  kinds: ResourceKind[]
}

export const inputs = reactive<Inputs>({
  diffResources: false,
  policy: PolicyTemplate,
  oldResource: '',
  resource: ResourceTemplate,
  context: ContextTemplate,
  config: ConfigTemplate,
  exceptions: '',
  clusterResources: '',
  imageData: ImageDataTemplate,
  customResourceDefinitions: CustomResourceDefinitionsTemplate
})

export const state = reactive<State>({
  kinds: []
})

export const reset = () => {
  inputs.diffResources = false
  inputs.policy = PolicyTemplate
  inputs.oldResource = ''
  inputs.exceptions = ''
  inputs.clusterResources = ''
  inputs.resource = ResourceTemplate
  inputs.context = ContextTemplate
  inputs.config = ConfigTemplate
  inputs.imageData = ImageDataTemplate
  inputs.customResourceDefinitions = CustomResourceDefinitionsTemplate
}

export const setDefaults = () => {
  init({
    policy: PolicyTemplate,
    diffResources: false,
    resource: ResourceTemplate,
    oldResource: '',
    context: ContextTemplate,
    config: ConfigTemplate,
    exceptions: '',
    clusterResources: '',
    imageData: ImageDataTemplate,
    customResourceDefinitions: CustomResourceDefinitionsTemplate
  })
}

export const init = (values: Inputs) => {
  const state = useState()

  if (typeof values.policy === 'string') {
    state.policy.value = values.policy
    inputs.policy = values.policy
  }

  if (typeof values.resource === 'string') {
    state.resource.value = values.resource
    inputs.resource = values.resource
  }

  if (values.diffResources !== undefined) {
    inputs.diffResources = values.diffResources
  }

  if (!values.oldResource && values.resource && inputs.diffResources) {
    state.oldResource.value = values.resource
    inputs.oldResource = values.resource
  }

  if (values.oldResource) {
    inputs.diffResources = true
    state.oldResource.value = values.oldResource
    inputs.oldResource = values.oldResource
  } else {
    inputs.diffResources = false
    state.oldResource.value = null
    inputs.oldResource = ''
  }

  if (typeof values.context === 'string') {
    state.context.value = values.context
    inputs.context = values.context
  }

  if (typeof values.config === 'string') {
    state.config.value = values.config
    inputs.config = values.config
  }

  if (typeof values.exceptions === 'string') {
    state.exceptions.value = values.exceptions
    inputs.exceptions = values.exceptions
  }

  if (typeof values.clusterResources === 'string') {
    state.clusterResources.value = values.clusterResources
    inputs.clusterResources = values.clusterResources
  }

  if (typeof values.customResourceDefinitions === 'string') {
    state.customResourceDefinitions.value = values.customResourceDefinitions
    inputs.customResourceDefinitions = values.customResourceDefinitions
  }

  if (typeof values.imageData === 'string') {
    state.imageData.value = values.imageData
    inputs.imageData = values.imageData
  }

  if (typeof values.vapBindings === 'string') {
    state.vapBindings.value = values.vapBindings
    inputs.vapBindings = values.vapBindings
  }

  state.name.value = values.name || ''
}

export const update = (values: Inputs) => {
  if (values.policy) {
    inputs.policy = values.policy
  }

  if (values.oldResource) {
    inputs.oldResource = values.oldResource
  }

  if (values.resource) {
    inputs.resource = values.resource
  }

  if (values.context) {
    inputs.context = values.context
  }

  if (values.config) {
    inputs.config = values.config
  }

  if (values.exceptions) {
    inputs.exceptions = values.exceptions
  }

  if (values.clusterResources) {
    inputs.clusterResources = values.clusterResources
  }

  if (values.customResourceDefinitions) {
    inputs.customResourceDefinitions = values.customResourceDefinitions
  }

  if (values.imageData) {
    inputs.imageData = values.imageData
  }
}

export const populate = () => {
  const state = useState()

  if (state.policy.value) {
    inputs.policy = state.policy.value
  }

  if (state.oldResource.value) {
    inputs.diffResources = true
    inputs.oldResource = state.oldResource.value
  }

  if (state.resource.value) {
    inputs.resource = state.resource.value
  }

  if (state.context.value) {
    inputs.context = state.context.value
  }

  if (state.config.value) {
    inputs.config = state.config.value
  }

  if (state.exceptions.value) {
    inputs.exceptions = state.exceptions.value
  }

  if (state.clusterResources.value) {
    inputs.clusterResources = state.clusterResources.value
  }

  if (state.customResourceDefinitions.value) {
    inputs.customResourceDefinitions = state.customResourceDefinitions.value
  }

  if (state.imageData.value) {
    inputs.imageData = state.imageData.value
  }
}
