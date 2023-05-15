export const resolveAPI = () => {
    const api: string = import.meta.env.VITE_API_HOST || `${window.location.protocol}//${window.location.origin}/${window.location.pathname}`

    return api.trim().replace(/\/$/, '')
}