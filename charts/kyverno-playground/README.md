# kyverno-playground

Kyverno Playground Web Application

![Version: 0.5.3](https://img.shields.io/badge/Version-0.5.3-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v0.5.3](https://img.shields.io/badge/AppVersion-v0.5.3-informational?style=flat-square)

## About

The Kyverno Playground is a web service that simulates [Kyverno](https://github.com/kyverno/kyverno) behaviuor, you can experiment and play with Kyverno policies directly in your browser

The service receives a configuration, resource and policies definitions, runs the Kyverno engine, and returns the results of evaluating policies against resources.

The playground currently supports:
- Validation rules
- Mutation rules
- Image verification rules
- Generate rules

**NOTES:**
- This tool only works with public image registries
- No data is gathered, stored, or shared

## Features

The playground frontend offers a rich feature set:
- Supports admission information like `username`, `groups`, `roles` and `cluster roles`
- Saving and loading state from the local storage
- Loading policies and resources from the Kyverno catalog
- Sharing state with simple links
- Comes with a tutorial to learn Kyverno easily

## Installing the Chart

Add `kyverno-playground` Helm repository:

```shell
helm repo add kyverno-playground https://kyverno.github.io/playground/
```

Install `kyverno-playground` Helm chart:

```shell
helm install kyverno-playground --namespace kyverno --create-namespace kyverno-playground/kyverno-playground
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| nameOverride | string | `""` | Name override |
| fullnameOverride | string | `""` | Full name override |
| replicaCount | int | `1` | Number of pod replicas |
| sponsor | string | `""` | Optional sponsor text |
| tufRootMountPath | string | `"/.sigstore"` | A writable volume to use for the TUF root initialization. |
| image.registry | string | `"ghcr.io"` | Image registry |
| image.repository | string | `"kyverno/playground"` | Image repository |
| image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| image.tag | string | `nil` | Image tag (will default to app version if not set) |
| imagePullSecrets | list | `[]` | Image pull secrets |
| priorityClassName | string | `""` | Priority class name |
| serviceAccount.create | bool | `true` | Create service account |
| serviceAccount.annotations | object | `{}` | Service account annotations |
| serviceAccount.name | string | `""` | Service account name (required if `serviceAccount.create` is `false`) |
| podAnnotations | object | `{}` | Pod annotations |
| podSecurityContext | object | `{"fsGroup":2000}` | Pod security context |
| securityContext | object | See [values.yaml](values.yaml) | Container security context |
| service.type | string | `"ClusterIP"` | Service type |
| service.port | int | `8080` | Service port |
| livenessProbe | object | `{"httpGet":{"path":"/","port":"http"}}` | Liveness probe |
| readinessProbe | object | `{"httpGet":{"path":"/","port":"http"}}` | Readiness probe |
| ingress.enabled | bool | `false` | Enable ingress |
| ingress.className | string | `""` | Ingress class name |
| ingress.annotations | object | `{}` | Ingress annotations |
| ingress.hosts | list | `[]` | Ingress hosts |
| ingress.tls | list | `[]` | Ingress tls |
| resources.limits | string | `nil` | Container resource limits |
| resources.requests | string | `nil` | Container resource requests |
| autoscaling.enabled | bool | `false` | Enable autoscaling |
| autoscaling.minReplicas | int | `1` | Min number of replicas |
| autoscaling.maxReplicas | int | `100` | Max number of replicas |
| autoscaling.targetCPUUtilizationPercentage | int | `80` | Target CPU utilisation |
| autoscaling.targetMemoryUtilizationPercentage | string | `nil` | Target Memory utilisation |
| nodeSelector | object | `{}` | Node selector |
| tolerations | list | `[]` | Tolerations |
| affinity | object | `{}` | Affinity |
| clusterRoles | list | `[]` | Cluster roles |
| roles | list | `[]` | Cluster roles |
| extraArgs | object | `{}` | Additonal container arguments |
| config.gin.mode | string | `"release"` | Gin mode (`release` or `debug`) |
| config.gin.cors | bool | `false` | Gin cors middleware |
| config.gin.logger | bool | `false` | Gin logger middleware |
| config.gin.maxBodySize | int | `2097152` | Gin max body size |
| config.server.host | string | `"0.0.0.0"` | Server host |
| config.server.port | int | `8080` | Server port |
| config.cluster.enabled | bool | `false` | Enable connected cluster mode |
| config.ui.sponsor | string | `""` | Sponsor name |
| config.engine.builtinCrds | list | `[]` | Builtin CRDs enabled (`argocd`, `cert-manager`, `prometheus-operator`, `tekton-pipelines`) |
| config.engine.localCrds | list | `[]` | Paths to folders containing yaml definitions for CRDs |
| config.versions | list | `[]` | list of additional Kyverno Playground versions |

## Source Code

* <https://github.com/kyverno/playground>

## Requirements

Kubernetes: `>=1.16.0-0`

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Nirmata |  | <https://kyverno.io/> |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.11.0](https://github.com/norwoodj/helm-docs/releases/v1.11.0)
