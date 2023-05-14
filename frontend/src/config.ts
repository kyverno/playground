import { useLocalStorage, usePreferredDark } from "@vueuse/core";
import { watch, computed } from "vue";

export type Config = {
    editorThemes: { name: string; theme: string; }[];
    layoutThemes: string[];
    onboarding: { text: string };
}

export type Policy = {
    url?: string;
    contextPath?: string;
    path: string;
    title: string;
}

type Example = {
    name: string;
    color?: string;
    url: string;
    subgroups?: { name: string; policies: Policy[]; url?: string; }[]
    policies?: Policy[] | string[];
}

const isDark = usePreferredDark()
export const layoutTheme = useLocalStorage<'light' | 'dark'>('config:layoutTheme', isDark.value ? 'dark' : 'light')
watch(isDark, (dark: boolean) => {
    layoutTheme.value = dark ? 'dark' : 'light'
})

export const btnColor = computed(() => {
    if (layoutTheme.value === 'dark') return 'secondary'

    return 'primary'
})

export const editorTheme = useLocalStorage('config:editorTheme', 'vs-dark')
export const hideNoMatch = useLocalStorage('config:hideNoMatch', false)
export const showOnboarding = useLocalStorage("onboarding:open", true)

export const options = {
    onboarding: {
        text: 'Notice: This tool only works with public image registries. No data is gathered, stored, or shared.',
    },
    layoutThemes: ['light', 'dark'],
    editorThemes: [
        { name: 'VS Dark', theme: 'vs-dark' }, 
        { name: 'VS Light', theme: 'vs' },
        { name: 'HC Black', theme: 'hc-black' }, 
        { name: 'HC Light', theme: 'hc-light' }
    ],
    examples: [
        {
            name: "Tutorials",
            url: 'https://raw.githubusercontent.com/kyverno/policies/main/',
            color: 'warning',
            subgroups: [
                {
                    name: 'RuleTypes',
                    url: '/tutorials',
                    policies: [
                        { path: 'validate', title: 'Validation' },
                        { path: 'mutation', title: 'Mutation' },
                        { path: 'generation', title: 'Generation' },
                        { path: 'verify-images', title: 'Verify Images' },
                    ]
                },
                {
                    name: 'ConfigMap Context',
                    policies: [
                        { url: '/tutorials/policies', path: 'allowed-pod-priorities', title: 'Allowed Pod Priorities' },
                        { path: 'other/exclude-namespaces-dynamically', contextPath: '/tutorials/context', title: 'Exclude Namespaces Dynamically' },
                    ]
                },
                {
                    name: 'API Call Context',
                    policies: [
                        { path: 'other/restrict-pod-count-per-node', contextPath: '/tutorials/context', title: 'Restrict Pod Count per Node' },
                        { path: 'other/restrict-ingress-host', contextPath: '/tutorials/context', title: 'Unique Ingress Host' },
                        { path: 'other/require-netpol', contextPath: '/tutorials/context', title: 'Require NetworkPolicy' },
                    ]
                },
                {
                    name: 'UPDATE Operations',
                    policies: [
                        { path: 'other/allowed-label-changes', contextPath: '/tutorials/context', title: 'Allowed Label Changes' },
                        { path: 'other/block-updates-deletes', contextPath: '/tutorials/context', title: 'Block Updates and Deletes' },
                    ]
                },
                {
                    name: 'Subject Configuration',
                    policies: [
                        { path: 'other/check-serviceaccount', contextPath: '/tutorials/context', title: 'Check ServiceAccount' },
                    ]
                }
            ]
        },
        {
            name: "Best Practices",
            url: 'https://raw.githubusercontent.com/kyverno/policies/main/best-practices',
            color: undefined,
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
        {
            name: "Pod Security Baseline",
            url: 'https://raw.githubusercontent.com/kyverno/policies/main/pod-security/baseline',
            color: undefined,
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
        {
            name: "Pod Security Restricted",
            url: 'https://raw.githubusercontent.com/kyverno/policies/main/pod-security/restricted',
            color: undefined,
            policies: [
                "disallow-capabilities-strict",
                "disallow-privilege-escalation",
                "require-run-as-non-root-user",
                "require-run-as-nonroot",
                "restrict-seccomp-strict",
                "restrict-volume-types",
            ]
        },
        {
            name: "Pod Security Subrule",
            url: 'https://raw.githubusercontent.com/kyverno/policies/main/pod-security/subrule',
            color: undefined,
            policies: [
                "podsecurity-subrule-baseline",
            ]
        },
        {
            name: "Pod Security Subrule Restricted",
            url: 'https://raw.githubusercontent.com/kyverno/policies/main/pod-security/subrule/restricted',
            color: undefined,
            policies: [
                "restricted-exclude-capabilities",
                "restricted-exclude-seccomp",
                "restricted-latest",
            ]
        },
        {
            name: "Other",
            url: 'https://raw.githubusercontent.com/kyverno/policies/main/other',
            color: undefined,
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
    ] as Example[]
}

export const useConfig = () => ({
    editorTheme,
    layoutTheme,
    showOnboarding,
    options,
    hideNoMatch
})
