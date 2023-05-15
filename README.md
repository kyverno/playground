# Kyverno Playground

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

Additionaly we publish docker images at ghcr.io/kyverno/playground and an helm chart repository is available at https://kyverno.github.io/playground.

### Install with Helm

Add `kyverno-playground` Helm repository:

```shell
helm repo add kyverno-playground https://kyverno.github.io/playground/
```

Install `kyverno-playground` Helm chart:

```shell
helm install kyverno-playground --namespace kyverno --create-namespace kyverno-playground/kyverno-playground
```

## Screenshots

![Kyverno Playground - Layout](./frontend/screens/layout.png?raw=true)

<hr />

![Kyverno Playground - Examples](./frontend/screens/examples.png?raw=true)

<hr />

![Kyverno Playground - Validation Results](./frontend/screens/results.png?raw=true)

<hr />

![Kyverno Playground - DarkMode](./frontend/screens/darkmode.png?raw=true)
