import { useLocalStorage } from "@vueuse/core"
import { Ref, ref, watch } from "vue"

const persisted = useLocalStorage<string>('persist:list', '')

export const getPersisted = (): Ref<string[]> => {
    const list = ref<string[]>([])

    watch(persisted, (content: string) => {
        list.value = (content || '').split(';;').filter(l => !!l)
    }, { immediate: true })

    return list
}

export const createInput = (name: string) => {
    name = name.replaceAll(';;', ';').trim()
    const policy = useLocalStorage<string>(`persist:policy:${name}`, null)
    const resource = useLocalStorage<string>(`persist:resource:${name}`, null)
    const context = useLocalStorage<string>(`persist:context:${name}`, null)

    persisted.value = [...new Set([...getPersisted().value, name])].join(';;')

    return {
        policy,
        resource,
        context,
        name
    }
}

export const removeInput = (name: string) => {
    const input = createInput(name)

    input.policy.value = null
    input.resource.value = null
    input.context.value = null

    name = name.replaceAll(';;', ';').trim()
    const list = getPersisted()

    persisted.value = list.value.filter(l => l !== name).join(';;')
}

export const useLocalInput = (name: string) => {
    name = name.replaceAll(';;', ';').trim()
    const policy = useLocalStorage<string>(`persist:policy:${name}`, null)
    const resource = useLocalStorage<string>(`persist:resource:${name}`, null)
    const context = useLocalStorage<string>(`persist:context:${name}`, null)

    const list = getPersisted()

    persisted.value = [...new Set([...list.value, name])].join(';;')

    return {
        input: {
            policy,
            resource,
            context,
            name
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