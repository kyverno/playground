# Kyverno Playground

![release](https://github.com/kyverno/playground/workflows/release/badge.svg)
![ci](https://github.com/kyverno/playground/workflows/ci/badge.svg)
![image](https://github.com/kyverno/playground/workflows/image/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/kyverno/playground/backend)](https://goreportcard.com/report/github.com/kyverno/playground/backend)
![License: Apache-2.0](https://img.shields.io/github/license/kyverno/playground?color=blue)
[![GitHub Repo stars](https://img.shields.io/github/stars/kyverno/playground)](https://github.com/kyverno/playground/stargazers)

The **public version of the Playground** is available at https://playground.kyverno.io.

## About

The Kyverno Playground is a web service that simulates [Kyverno](https://github.com/kyverno/kyverno) behaviour, you can experiment and play with Kyverno policies directly in your browser

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

### Context and Variables

It is currently not possible to add variables from external resources or do actual API calls.

It is only possible to mock variables using the variables configuration in the context input.

### Multiple Manifests

It is supported to define multiple policies and/or resources as inputs.

Context and variables will be shared for all executions.

### Load Manifests

The "File" Button loads a local YAML file as input.

The "URL" Button loads a manifest from an external URL, example: https://raw.githubusercontent.com/kyverno/policies/main/best-practices/disallow-latest-tag/disallow-latest-tag.yaml

## Install

Kyverno Playground releases are available at https://github.com/kyverno/playground/releases.

Additionaly we publish docker images at [ghcr.io/kyverno/playground](https://github.com/kyverno/playground/pkgs/container/playground) and an helm chart repository is available at https://kyverno.github.io/playground.

### Install with Helm

Add `kyverno-playground` Helm repository:

```shell
helm repo add kyverno-playground https://kyverno.github.io/playground/
```

Install `kyverno-playground` Helm chart:

```shell
helm upgrade --install kyverno-playground --namespace kyverno --create-namespace --wait kyverno-playground/kyverno-playground
```

Install `kyverno-playground` Helm chart (without configuring an Helm repository):
```shell
helm upgrade --install kyverno-playground --namespace kyverno --create-namespace --wait --repo https://kyverno.github.io/playground kyverno-playground
```

Install `kyverno-playground` local Helm chart:
```shell
helm upgrade --install kyverno-playground --namespace kyverno --create-namespace --wait ./charts/kyverno-playground
```

### Install and run locally

Alternatively, you can install and run the Playground locally. This will allow you to connect the Playground to a real cluster.

Please read the [Cluster connected docs](./docs/CLUSTER.md).

## Custom resources

The Playground uses openapi schemas to load resources from yaml content. To load a resource correctly the Playground needs the corresponding openapi schema.

By default, all Kubernetes builtin resources are supported. To work with custom resources you need to provide the custom resource definition.

Providing custom resource definitions can be done in different ways:
* Using the `--builtin-crds` flag in the backend (see the list of [supported built-in CRDs](#supported-built-in-custom-resource-definitions))
* Using the `--local-crds` flag in the backend, pointing to a directory containing yaml CRD definitions
* Paste your CRD yaml definitions directly in the frontend

### Supported built-in custom resource definitions

The following CRDs are embedded in the Playground backend and can be enabled with the `--builtin-crds` flag:

| Name | Flag |
| --- | --- |
| ArgoCD | `--builtin-crds=argocd` |
| Cert Manager | `--builtin-crds=cert-manager` |
| Tekton Pipeline | `--builtin-crds=tekton-pipeline` |
| Prometheus Operator | `--builtin-crds=prometheus-operator` |

## Build

Instructions for building and running the Playground from source code is available in the [docs](./docs) section.

## Screenshots

![Kyverno Playground - Layout](./frontend/screens/layout.png?raw=true)

<hr />

![Kyverno Playground - Examples](./frontend/screens/examples.png?raw=true)

<hr />

![Kyverno Playground - Validation Results](./frontend/screens/results.png?raw=true)

<hr />

![Kyverno Playground - DarkMode](./frontend/screens/darkmode.png?raw=true)
