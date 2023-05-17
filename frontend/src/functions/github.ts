import { ContextTemplate, ResourceTemplate } from "@/assets/templates";
import { init } from "@/store";

const baseURL = 'https://raw.githubusercontent.com/kyverno/policies/main'

export type Policy = {
    url?: string;
    contextPath?: string;
    resourceFile?: string;
    path: string;
    title: string;
}

export type LoadConfig = {
    start?: () => void,
    success?: (values: [string, string, string]) => void
    error?: (err : Error) => void
    finished?: () => void
}

export const loadPolicy = async (url: string, policy: Policy, config?: LoadConfig) => {
    const folder = policy.path
    const contextPath = policy.contextPath
    const name = folder.split('/').pop()

    try {
        const contextURL = contextPath ? `${contextPath}/${name}.yaml` : `${url}/${folder}/context.yaml`
        const resourceFile = policy.resourceFile || 'resource.yaml'

        const promises = [
            fetch(`${url}/${folder}/${name}.yaml`).then((resp) => resp.text()),
            fetch(`${url}/${folder}/${resourceFile}`).then((resp) => {
                if (resp.status === 404) {
                    return fetch(`${url}/${folder}/resources.yaml`)
                }

                return resp
            }).then((resp) => {
                if (resp.status === 404) {
                    return ResourceTemplate
                }

                return resp.text()
            }),
            fetch(contextURL).then((resp) => {
                if (resp.status === 404) {
                    return ContextTemplate
                }

                return resp.text()
            }),
        ]

        config?.start?.();
        const values = await Promise.all(promises);

        config?.success?.(values as [string, string, string]);
    } catch (err) {
        config?.error?.(err as Error)
    } finally {
        config?.finished?.()
    }
};

export const loadFromRepo = (path: string, resource?: string, config?: LoadConfig) => {
    return loadPolicy(baseURL, {
        path: path.replace(/^\/+/, '').replace(/\/+$/, '').trim(),
        resourceFile: resource,
        title: ''
    }, {
        success([policy, resource, context]) {
            init({ policy, resource, context })
        },
        error(err) {
            console.error(err)
        },
        ...(config || {}),
    })
}