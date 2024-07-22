import { createInput } from '@/composables'
import { init, inputs } from '@/store'
import { parse, stringify } from 'yaml'

export type ProfileExport = {
  date: string
  version: string
  profiles?: {
    name?: string
    policies?: string
    resources?: string
    oldResources?: string
    context?: string
    config?: string
    customResourceDefinitions?: string
    clutserResources?: string
    exceptions?: string
    imageData?: string
    vapBindings?: string
  }[]
}

export const convertProfiles = (current: boolean, profiles: string[]): string => {
  const exports = []

  if (current) {
    exports.push({
      name: 'current-state',
      policies: inputs.policy,
      resources: inputs.resource,
      oldResources: inputs.oldResource,
      context: inputs.context,
      config: inputs.config,
      customResourceDefinitions: inputs.customResourceDefinitions,
      clusterResources: inputs.clusterResources,
      exceptions: inputs.exceptions,
      imageData: inputs.imageData,
      vapBindings: inputs.vapBindings
    })
  }

  profiles.map((p) => {
    const {
      policy,
      resource,
      context,
      config,
      oldResource,
      customResourceDefinitions,
      clusterResources,
      exceptions,
      imageData
    } = createInput(p)

    exports.push({
      name: p,
      policies: policy.value,
      resources: resource.value,
      oldResources: oldResource.value,
      context: context.value,
      config: config.value,
      customResourceDefinitions: customResourceDefinitions.value,
      clusterResources: clusterResources.value,
      exceptions: exceptions.value,
      imageData: imageData.value,
      vapBindings: inputs.vapBindings
    })
  })

  const state: ProfileExport = {
    date: new Date().toISOString().slice(0, 10),
    version: APP_VERSION,
    profiles: exports
  }

  return stringify(state, null, { lineWidth: 0 })
}

export const importProfiles = async (content: string) => {
  const state: ProfileExport | undefined = parse(content)
  if (!state || typeof state !== 'object') {
    throw new Error('invalid import file')
  }

  if (!Array.isArray(state.profiles)) {
    throw new Error('invalid import file')
  }

  if (!state.profiles.length) return

  if (state.profiles[0]?.name === 'current-state') {
    const currentState = state.profiles.shift()

    init({
      policy: currentState?.policies,
      resource: currentState?.resources,
      oldResource: currentState?.oldResources,
      context: currentState?.context,
      config: currentState?.config,
      customResourceDefinitions: currentState?.customResourceDefinitions,
      clusterResources: currentState?.clutserResources,
      exceptions: currentState?.exceptions,
      imageData: currentState?.imageData,
      vapBindings: currentState?.vapBindings
    })
  }

  state.profiles
    .filter((p) => !!p)
    .forEach((profile) => {
      if (!profile.name) {
        console.warn('invalid profile, no name found')
        return
      }

      createInput(profile.name, {
        policy: profile.policies,
        resource: profile.resources,
        oldResource: profile?.oldResources,
        context: profile.context,
        config: profile?.config,
        customResourceDefinitions: profile?.customResourceDefinitions,
        clusterResources: profile?.clutserResources,
        exceptions: profile?.exceptions,
        imageData: profile.imageData,
        vapBindings: profile?.vapBindings
      })
    })
}
