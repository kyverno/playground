import { useLocalStorage } from "@vueuse/core"
import { Ref, ref, watch } from "vue"

export type Inputs = {
    name?: string;
    diffResources?: boolean;
    policy?: string | null;
    oldResource?: string | null;
    resource?: string | null;
    context?: string | null;
    config?: string | null;
    customResourceDefinitions?: string | null;
}

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
    const oldResource = useLocalStorage<string | null>(`persist:resource:old:${name}`, defaults?.oldResource || null)
    const context = useLocalStorage<string | null>(`persist:context:${name}`, defaults?.context || null)
    const config = useLocalStorage<string | null>(`persist:config:${name}`, defaults?.config || null)
    const customResourceDefinitions = useLocalStorage<string | null>(`persist:crds:${name}`, defaults?.customResourceDefinitions || null)

    persisted.value = [...new Set([...getPersisted().value, name])].join(';;')

    return {
        policy,
        resource,
        oldResource,
        context,
        config,
        customResourceDefinitions,
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
    if (input.config.value !== values.config) {
        input.config.value = values.config
    }
    if (input.customResourceDefinitions.value !== values.customResourceDefinitions) {
        input.customResourceDefinitions.value = values.customResourceDefinitions
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
    input.customResourceDefinitions.value = null

    name = name.replaceAll(';;', ';').trim()
    const list = getPersisted()

    persisted.value = list.value.filter(l => l !== name).join(';;')
}

export const useLocalInput = (name: string) => {
    name = name.replaceAll(';;', ';').trim()
    const policy = useLocalStorage<string>(`persist:policy:${name}`, null)
    const resource = useLocalStorage<string>(`persist:resource:${name}`, null)
    const oldResource = useLocalStorage<string>(`ppersist:resource:old:${name}`, null)
    const context = useLocalStorage<string>(`persist:context:${name}`, null)
    const config = useLocalStorage<string | null>(`persist:config:${name}`, null)
    const customResourceDefinitions = useLocalStorage<string | null>(`persist:crds:${name}`, null)

    const list = getPersisted()

    persisted.value = [...new Set([...list.value, name])].join(';;')

    return {
        input: {
            policy,
            resource,
            oldResource,
            context,
            config,
            customResourceDefinitions,
            name,
        },
        remove: () => {
            policy.value = null
            resource.value = null
            oldResource.value = null
            context.value = null
            config.value = null
            customResourceDefinitions.value = null

          persisted.value = list.value.filter(l => l !== name).join(';;')
        },
        list
    }
}
