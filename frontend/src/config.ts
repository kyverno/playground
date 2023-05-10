import { useLocalStorage, usePreferredDark } from "@vueuse/core";
import { reactive, watch } from "vue";

export type Config = {
    editorThemes: { name: string; theme: string; }[];
    examples: { [group: string]: { path: string; policies: string[] } }[]
}

const isDark = usePreferredDark()
export const layoutTheme = useLocalStorage<'light' | 'dark'>('config:layoutTheme', isDark.value ? 'dark' : 'light')
watch(isDark, (dark: boolean) => {
    layoutTheme.value = dark ? 'dark' : 'light'
})

export const editorTheme = useLocalStorage('config:editorTheme', 'vs-dark')

export const hideNoMatch = useLocalStorage('config:hideNoMatch', false)


export const policyLS = useLocalStorage<string>('persist:policy', null)
export const resourceLS = useLocalStorage<string>('persist:resource', null)
export const contextLS = useLocalStorage<string>('persist:context', null)

export const options = reactive({
    layoutThemes: ['light', 'dark'],
    editorThemes: [
        { name: 'VS Dark', theme: 'vs-dark' }, 
        { name: 'VS Light', theme: 'vs' },
        { name: 'HC Black', theme: 'hc-black' }, 
        { name: 'HC Light', theme: 'hc-light' }
    ],
    examples: {
        "Best Practices": {
            path: 'https://raw.githubusercontent.com/kyverno/policies/main/best-practices',
            policies: [
                "disallow-cri-sock-mount",
                "disallow-default-namespace",
                "disallow-empty-ingress-host",
                "disallow-helm-tiller",
                "disallow-latest-tag",
                "require-drop-all",
                "require-drop-cap-net-raw",
                "require-labels",
                "require-pod-requests-limits",
                "require-probes",
                "require-ro-rootfs",
                "restrict-image-registries",
                "restrict-node-port",
                "restrict-service-external-ips",
            ]
        },
        "Pod Security Baseline": {
            path: 'https://raw.githubusercontent.com/kyverno/policies/main/pod-security/baseline',
            policies: [
                "disallow-capabilities",
                "disallow-host-namespaces",
                "disallow-host-path",
                "disallow-host-ports-range",
                "disallow-host-ports",
                "disallow-host-process",
                "disallow-privileged-containers",
                "disallow-proc-mount",
                "disallow-selinux",
                "restrict-apparmor-profiles",
                "restrict-seccomp",
                "restrict-sysctls",
            ]
        },
        "Pod Security Restricted": {
            path: 'https://raw.githubusercontent.com/kyverno/policies/main/pod-security/restricted',
            policies: [
                "disallow-capabilities-strict",
                "disallow-privilege-escalation",
                "require-run-as-non-root-user",
                "require-run-as-nonroot",
                "restrict-seccomp-strict",
                "restrict-volume-types",
            ]
        },
        "Pod Security Subrule": {
            path: 'https://raw.githubusercontent.com/kyverno/policies/main/pod-security/subrule',
            policies: [
                "podsecurity-subrule-baseline",
            ]
        },
        "Pod Security Subrule Restricted": {
            path: 'https://raw.githubusercontent.com/kyverno/policies/main/pod-security/subrule/restricted',
            policies: [
                "restricted-exclude-capabilities",
                "restricted-exclude-seccomp",
                "restricted-latest",
            ]
        },
        "Cert Manager": {
            path: 'https://raw.githubusercontent.com/kyverno/policies/main/cert-manager',
            policies: [
                "limit-dnsnames",
                "limit-duration",
                "restrict-issuer",
            ]
        },
        "NGINX Ingress": {
            path: 'https://raw.githubusercontent.com/kyverno/policies/main/nginx-ingress',
            policies: [
                "disallow-ingress-nginx-custom-snippets",
                "restrict-annotations",
                "restrict-ingress-paths",
            ]
        },
        "Velero": {
            path: 'https://raw.githubusercontent.com/kyverno/policies/main/velero',
            policies: [
                "block-velero-restore",
                "validate-cron-schedule",
                "backup-all-volumes",
            ]
        },
        "Other": {
            path: 'https://raw.githubusercontent.com/kyverno/policies/main/other',
            policies: [
                "add-certificates-volume",
                "add-default-resources",
                "add-labels",
                "allowed-annotations",
                "allowed-pod-priorities",
                "check-env-vars",
                "require-base-image",
            ]
        },
    }
})

export const useConfig = () => ({
    editorTheme,
    layoutTheme,
    options,
    hideNoMatch,
    policy: policyLS,
    resource: resourceLS,
    context: contextLS
})
