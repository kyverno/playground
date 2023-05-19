import { createInput } from "@/composables"
import { init, inputs } from "@/store"
import { parse, stringify } from "yaml"

export type ProfileExport = {
    date: string;
    version: string;
    config?: string;
    profiles?: {
        name?: string;
        policies?: string;
        resources?: string;
        context?: string;
    }[]
}

export const convertProfiles = (current: boolean, profiles: string[], config: boolean): string => {
    const exports = []

    if (current) {
        exports.push({
            name: 'current-state',
            policies: inputs.policy,
            resources: inputs.resource,
            context: inputs.context,
        })
    }

    profiles.map(p => {
        const { policy, resource, context } = createInput(p)

        exports.push({
            name: p,
            policies: policy.value,
            resources: resource.value,
            context: context.value
        })
    })

    const state: ProfileExport = {
        date: (new Date()).toISOString().slice(0, 10),
        version: APP_VERSION,
        profiles: exports
    }

    if (config) {
        state.config = inputs.config
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

    if (!state.profiles.length) return;

    if (state.profiles[0]?.name === 'current-state') {
        const currentState = state.profiles.shift()

        init({
            policy: currentState?.policies,
            resource: currentState?.resources,
            context: currentState?.context,
        })
    }

    if (state.config) {
        inputs.config = state.config
    }

    state.profiles.filter((p => !!p)).forEach((profile) => {
        if (!profile.name) {
            console.warn('invalid profile, no name found')
            return;
        }

        createInput(profile.name, {
            policy: profile.policies,
            resource: profile.resources,
            context: profile.context
        })
    })
} 