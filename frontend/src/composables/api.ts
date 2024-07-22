import { resolveAPI } from '@/utils'
import { type Ref, reactive, ref } from 'vue'
import { stringify } from 'yaml'

export type Config = {
  sponsor: string
  cluster: boolean
}

export type ListRequest = {
  apiVersion: string
  kind: string
  namespace: string
}

export type ResourceRequest = {
  apiVersion: string
  kind: string
  namespace: string
  name: string
}

export type ResourceKind = {
  apiVersion: string
  kind: string
  clusterScoped: boolean
}

export type Resource = ResourceKind & { title: string }

export type SearchResult = { namespace: string; name: string }

export const config = reactive<Config>({
  sponsor: '',
  cluster: false
})

const fetchConfig = (api: string, loading: Ref<boolean>, error: Ref<Error | undefined>) => () => {
  loading.value = true
  error.value = undefined

  fetch(`${api}/config`, {
    method: 'GET',
    mode: 'cors',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'application/json'
    }
  })
    .then(async (resp) => {
      if (resp.status > 200) {
        throw new Error(await resp.text())
      }

      return resp.json()
    })
    .then(({ sponsor, cluster }: Config) => {
      config.sponsor = sponsor
      config.cluster = cluster
    })
    .catch((err) => {
      error.value = err
    })
    .finally(() => {
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

const fetchWrapper = <T, R = undefined>(method: string, api: string, request?: R): Promise<T> => {
  return fetch(api, {
    method: method,
    mode: 'cors',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'application/json'
    },
    ...(request ? { body: JSON.stringify(request) } : {})
  }).then<T>(async (resp) => {
    if (resp.status > 200) {
      throw new Error(await resp.text())
    }

    return resp.json()
  })
}

const fetchNamespaces = (api: string) => () => fetchWrapper<string[]>('GET', `${api}/cluster/namespaces`)

const fetchResources = (api: string) => (request: ListRequest) => {
  return fetchWrapper<SearchResult[]>('GET', `${api}/cluster/search?` + new URLSearchParams(request).toString())
}

const fetchResource = (api: string) => (request: ResourceRequest) => {
  return fetchWrapper<object>('GET', `${api}/cluster/resource?` + new URLSearchParams(request).toString())
}

const fetchKinds = (api: string) => () => fetchWrapper<ResourceKind[]>('GET', `${api}/cluster/kinds`)

export const useAPI = <T>() => {
  const api = resolveAPI()

  return {
    error: ref<Error | undefined>(),
    loading: ref<boolean>(false),
    data: ref<T>(),
    namespaces: fetchNamespaces(api),
    resources: fetchResources(api),
    resource: fetchResource(api),
    kinds: fetchKinds(api)
  }
}

export const resourceToYAML = (resource: object): string => stringify(resource, { lineWidth: 0 })

export const resourcesToYAML = (resources: object[]): string => resources.map((r) => resourceToYAML(r)).join('---\n')
