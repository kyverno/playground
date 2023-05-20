import { resolveAPI } from "@/utils";
import { Ref, reactive, ref } from "vue";
import { stringify } from "yaml";

export type Config = {
    sponsor: string;
    cluster: boolean;
}

export type ListRequest = {
    apiVersion: string;
    kind: string;
    namespace?: string;
    selector?: { [key: string]: string };
}

export type ResourceRequest = {
    apiVersion: string;
    kind: string;
    namespace?: string;
    name: string;
}

export const config = reactive<Config>({
    sponsor: '',
    cluster: false,
})

const fetchConfig = (api: string, loading: Ref<boolean>, error: Ref<Error | undefined>) => () => {
    loading.value = true
    error.value = undefined

    fetch(`${api}/config`, {
        method: "POST",
        mode: "cors",
        cache: "no-cache",
        headers: {
            "Content-Type": "application/json",
        },
    }).then(async (resp) => {
        if (resp.status > 200) {
            throw new Error(await resp.text());
        }

        return resp.json()
    }).then(({ sponsor, cluster}: Config) => {
        config.sponsor = sponsor
        config.cluster = cluster
    }).catch((err) => {
        error.value = err
    }).finally(() => {
        loading.value = false
    })

    return { error, loading, config }
}

export const useAPIConfig = () => {
    const api = resolveAPI()

    const loading = ref(false)
    const error = ref()
    return {
        fetch: fetchConfig(api, loading, error),
        config
    }
}

const fetchWrapper = <T, R = undefined>(api: string, request: R): Promise<T> => {
    return fetch(api, {
        method: "POST",
        mode: "cors",
        cache: "no-cache",
        headers: {
            "Content-Type": "application/json",
        },
        ...(request ? { body: JSON.stringify(request) } : {})
    }).then<T>(async (resp) => {
        if (resp.status > 200) {
            throw new Error(await resp.text());
        }

        return resp.json()
    })
}

const fetchNamespaces = (api: string) => () => fetchWrapper<string[]>(`${api}/namespaces`, undefined)

const fetchResources = (api: string) => (request: ListRequest) => fetchWrapper<{ namespace?: string, name: string }[], ListRequest>(`${api}/resources`, request)

const fetchResource = (api: string) => (request: ResourceRequest) => fetchWrapper<object, ResourceRequest>(`${api}/resource`, request)

export const useAPI = <T>() => {
    const api = resolveAPI()

    return {
        error: ref<Error | undefined>(),
        loading: ref<boolean>(false),
        data: ref<T>(),
        namespaces: fetchNamespaces(api),
        resources: fetchResources(api),
        resource: fetchResource(api)
    }
}

export const resourceToYAML = (resource: object): string => stringify(resource, { lineWidth: 0 })

export const resourcesToYAML = (resources: object[]): string => resources.map(r => resourceToYAML(r)).join("---\n")