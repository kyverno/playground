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
    panels: {
        policyInfo: 'Kyverno Policy Resource',
        resourceInfo: 'Kubernetes resources to apply the policies on',
        contextInfo: 'Context information like operation conext, variables and kubernetes version',
    },
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
            color: 'orange-darken-3',
            subgroups: [
                {
                    name: 'RuleTypes',
                    url: `tutorials/policies`,
                    policies: [
                        { path: 'validate', title: 'Validate Pod Labels' },
                        { path: 'mutation', title: 'Mutate Pod Annotations' },
                        { path: 'generation', title: 'Generate Quotas' },
                        { path: 'verify-images', title: 'Verify Image Signatures' },
                    ]
                },
                {
                    name: 'ConfigMap Context',
                    policies: [
                        { url: 'tutorials/policies', path: 'allowed-pod-priorities', title: 'Allowed Pod Priorities' },
                        { path: 'other/exclude-namespaces-dynamically', contextPath: 'tutorials/context', title: 'Exclude Namespaces Dynamically' },
                    ]
                },
                {
                    name: 'API Call Context',
                    policies: [
                        { path: 'other/restrict-pod-count-per-node', contextPath: 'tutorials/context', title: 'Restrict Pod Count per Node' },
                        { path: 'other/restrict-ingress-host', contextPath: 'tutorials/context', title: 'Unique Ingress Host' },
                        { path: 'other/require-netpol', contextPath: 'tutorials/context', title: 'Require NetworkPolicy' },
                    ]
                },
                {
                    name: 'UPDATE Operations',
                    policies: [
                        { path: 'other/allowed-label-changes', contextPath: 'tutorials/context', title: 'Allowed Label Changes' },
                        { path: 'other/block-updates-deletes', contextPath: 'tutorials/context', title: 'Block Updates and Deletes' },
                    ]
                },
                {
                    name: 'Subject Configuration',
                    policies: [
                        { path: 'other/check-serviceaccount', contextPath: 'tutorials/context', title: 'Check ServiceAccount' },
                    ]
                }
            ]
        },
        {
            name: "Best Practices",
            url: 'https://raw.githubusercontent.com/kyverno/policies/main/best-practices',
            color: undefined,
            policies: [
                { path: 'disallow-cri-sock-mount', title: 'Disallow CRI socket mounts' },
                { path: 'disallow-default-namespace', title: 'Disallow Default Namespace' },
                { path: 'disallow-empty-ingress-host', title: 'Disallow empty Ingress host' },
                { path: 'disallow-cri-sock-mount', title: 'Disallow CRI socket mounts' },
                { path: 'disallow-latest-tag', title: 'Disallow Latest Tag' },
                { path: 'require-drop-all', title: 'Drop All Capabilities' },
                { path: 'require-drop-cap-net-raw', title: 'Drop CAP_NET_RAW' },
                { path: 'require-labels', title: 'Require Labels' },
                { path: 'require-pod-requests-limits', title: 'Require Limits and Requests' },
                { path: 'require-probes', title: 'Require Pod Probes' },
                { path: 'restrict-image-registries', title: 'Restrict Image Registries' },
            ]
        },
        {
            name: "Pod Security",
            color: undefined,
            subgroups: [
                {
                    name: 'Baseline',
                    url: 'https://raw.githubusercontent.com/kyverno/policies/main/pod-security/baseline',
                    policies: [
                        { path: "disallow-capabilities", title: "Disallow Capabilities" },
                        { path: "disallow-host-namespaces", title: "Disallow  Host Namespaces" },
                        { path: "disallow-host-path", title: "Disallow hostPath" },
                        { path: "disallow-host-ports-range", title: "Disallow hostPorts Range (Alternate)" },
                        { path: "disallow-host-ports", title: "Disallow hostPorts" },
                        { path: "disallow-host-process", title: "Disallow hostProcess" },
                        { path: "disallow-privileged-containers", title: "Disallow Privileged Containers" },
                        { path: "disallow-proc-mount", title: "Disallow procMount" },
                        { path: "disallow-selinux", title: "Disallow SELinux" },
                        { path: "restrict-apparmor-profiles", title: "Restrict AppArmor" },
                        { path: "restrict-seccomp", title: "Restrict Seccomp" },
                        { path: "restrict-sysctls", title: "Restrict sysctls" },
                    ]
                },
                {
                    name: 'Restricted',
                    url: 'https://raw.githubusercontent.com/kyverno/policies/main/pod-security/restricted',
                    policies: [
                        { path: "disallow-capabilities-strict", title: "Disallow Capabilities (Strict)" },
                        { path: "disallow-privilege-escalation", title: "Disallow Privilege Escalation" },
                        { path: "require-run-as-non-root-user", title: "Require Run As Non-Root User" },
                        { path: "require-run-as-nonroot", title: "Require runAsNonRoot" },
                        { path: "restrict-seccomp-strict", title: "Restrict Seccomp (Strict)" },
                        { path: "restrict-volume-types", title: "Restrict Volume Types" },
                    ]
                },
                {
                    name: 'Subrule',
                    url: 'https://raw.githubusercontent.com/kyverno/policies/main/pod-security/subrule',
                    policies: [
                        { path: "podsecurity-subrule-baseline", title: "Baseline Pod Security Standards" },
                    ]
                },
                {
                    name: 'Subrule Restricted',
                    url: 'https://raw.githubusercontent.com/kyverno/policies/main/pod-security/subrule/restricted',
                    policies: [
                        { path: "restricted-exclude-capabilities", title: "Restricted Exclude Capabilities" },
                        { path: "estricted-exclude-seccomp", title: "Restricted Exclude SECComp" },
                        { path: "restricted-latest", title: "Restricted Pod Security Standards" },
                    ]
                },
            ],
        },
        {
            name: "Other",
            url: 'https://raw.githubusercontent.com/kyverno/policies/main/other',
            color: undefined,
            policies: [
                { path: "add-certificates-volume", title: "Add Certificates as a Volume" },
                { path: "add-default-resources", title: "Add Default Resources" },
                { path: "add-labels", title: "Add Labels" },
                { path: "allowed-annotations", title: "Allowed Annotations" },
                { path: "check-env-vars", title: "Check Environment Variables" },
                { path: "require-base-image", title: "Check Image Base" },
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
