import { useLocalStorage } from "@vueuse/core"
import { Ref, ref, watch } from "vue"

export type Inputs = { name?: string, policy?: string | null, resource?: string | null; context?: string | null; config?: string | null }

const persisted = useLocalStorage<string>('persist:list', '')

export const getPersisted = (): Ref<string[]> => {
    const list = ref<string[]>([])

    watch(persisted, (content: string) => {
        list.value = (content || '').split(';;').filter(l => !!l)
    }, { immediate: true })

    return list
}

export const createInput = (name: string, defaults?: Inputs) => {
    name = name.replaceAll(';;', ';').trim()
    const policy = useLocalStorage<string | null>(`persist:policy:${name}`, defaults?.policy || null)
    const resource = useLocalStorage<string | null>(`persist:resource:${name}`, defaults?.resource || null)
    const context = useLocalStorage<string | null>(`persist:context:${name}`, defaults?.context || null)
    const config = useLocalStorage<string | null>(`persist:config:${name}`, defaults?.context || null)

    persisted.value = [...new Set([...getPersisted().value, name])].join(';;')

    return {
        policy,
        resource,
        context,
        config,
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
    if (input.context.value !== values.context) {
        input.context.value = values.context
    }
    if (input.config.value !== values.config) {
        input.config.value = values.config
    }

    return input
}

export const removeInput = (name: string) => {
    const input = createInput(name)

    input.policy.value = null
    input.resource.value = null
    input.context.value = null
    input.config.value = null

    name = name.replaceAll(';;', ';').trim()
    const list = getPersisted()

    persisted.value = list.value.filter(l => l !== name).join(';;')
}

export const useLocalInput = (name: string) => {
    name = name.replaceAll(';;', ';').trim()
    const policy = useLocalStorage<string>(`persist:policy:${name}`, null)
    const resource = useLocalStorage<string>(`persist:resource:${name}`, null)
    const context = useLocalStorage<string>(`persist:context:${name}`, null)
    const config = useLocalStorage<string | null>(`persist:config:${name}`, null)

    const list = getPersisted()

    persisted.value = [...new Set([...list.value, name])].join(';;')

    return {
        input: {
            policy,
            resource,
            context,
            config,
            name,
        },
        remove: () => {
            policy.value = null
            resource.value = null
            context.value = null

            persisted.value = list.value.filter(l => l !== name).join(';;')
        },
        list
    }
}