import { ContextTemplate, PolicyTemplate, ResourceTemplate } from "@/assets/templates";
import { useLocalStorage } from "@vueuse/core";

export const loadedPolicy = useLocalStorage<string>('loaded:policy', PolicyTemplate);
export const loadedContext = useLocalStorage<string>('loaded:context', ContextTemplate);
export const loadedResource = useLocalStorage<string>('loaded:resource', ResourceTemplate);
export const loadedState = useLocalStorage<string>('loaded:state', '')

export const useState = () => ({
    policy: loadedPolicy,
    resource: loadedResource,
    context: loadedContext,
    name: loadedState
})