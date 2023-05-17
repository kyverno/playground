import { ContextTemplate, PolicyTemplate, ResourceTemplate } from "@/assets/templates";
import { useLocalStorage } from "@vueuse/core";

export const loadedPolicy = useLocalStorage<string>('loaded:policy', PolicyTemplate);
export const loadedContext = useLocalStorage<string>('loaded:context', ContextTemplate);
export const loadedResource = useLocalStorage<string>('loaded:resource', ResourceTemplate);
export const loadedState = useLocalStorage<string>('loaded:state', '')

export type State = { name: string; policy?: string, resource?: string; context?: string }

const update = (values: State) => {
    if (values.policy) {
        loadedPolicy.value = values.policy;
    }

    if (values.resource) {
        loadedResource.value = values.resource;
    }

    if (values.context) {
        loadedContext.value = values.context;
    }
    
    loadedState.value = values.name;
}

export const useState = () => ({
    policy: loadedPolicy,
    resource: loadedResource,
    context: loadedContext,
    name: loadedState,
    update
})