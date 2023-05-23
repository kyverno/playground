# Cluster connected docs

This docs contains informations related to running the Playground connected to a real cluster

## Abstract

While the public [Kyverno Playground](https://playground.kyverno.io) is not connected to a Kubernetes cluster, it is possible to run it locally and provide a kubeconfig file to the backend to run the Playground in cluster connected mode.

## Features

When the Playground backend is connected to a cluster, the following features are enabled:
- Load policies and/or resources directly from the cluster
- Evaluate generate policies with `clone` statements
- Use policies with context entries referencing real cluster resources

## Warning

Kyverno policies can create or update resources in a cluster.
Using a kubeconfig with read-only access should limit the risk of accidentally modifying existing resources.

## Install

You can install the Playground by downloading the prebuilt binary from the [GitHub release page](https://github.com/kyverno/playground/releases).

An example install script can be found below:

```shell
VERSION="latest"
OS="Darwin"       # Darwin, Linux, or Windows
ARCH="x86_64"     # x86_64, arm64, or i386

get_latest_release() {
    curl --silent "https://api.github.com/repos/kyverno/playground/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

if [ $VERSION = "latest" ]
then
    VERSION=$(get_latest_release)
fi

FILE_NAME="playground_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/kyverno/playground/releases/download/${VERSION}/${FILE_NAME}"

curl -LO $DOWNLOAD_URL
tar -xvf $FILE_NAME
```

The Playground binary is named `kyverno-playground`.

## Run in cluster connected mode

To run the Playground in cluster connected mode, you need to provide the `--kubeconfig` flag and point it to your local kubeconfig.

```shell
./kyverno-playground --kubeconfig=~/.kube/config
```

## Run with docker

If you want to run the Playground using Docker, run:

```shell
docker run --rm \
    -p 8080:8080 \
    -v ~/.kube/config:/.kube/config:ro \
    ghcr.io/kyverno/playground --kubeconfig=/.kube/config
```

**Warning:** Depending on the OS, docker may run in a virtual machine and running the Playground may not be able to access a private network.

## Accessing the Plaground

With default flags, the Playground will be available at http://localhost:8080.
