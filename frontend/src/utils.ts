export const resolveAPI = () => {
    const api: string = import.meta.env.VITE_API_HOST || `${window.location.origin}${window.location.pathname}`

    return api.trim().replace(/\/+$/, '')
}

export const mergeResources = (a: string, b: string): string => {
    return `${a.trim()}
    ---
    ${b.trim()}`
}