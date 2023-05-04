import { reactive } from "vue";

export type Config = {
    theme: 'light' | 'dark'
    examples: { [group: string]: { path: string; policies: string[] } }[]
}

export const config = reactive({
    theme: 'light',
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
    }
})