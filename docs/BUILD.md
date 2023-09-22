# Build docs

This docs contains informations related to building the Playground

## Required

* Requires Go >= v1.21
* Requires NodeJS >= v18.0

## Makefile

We provide and maintain a [Makefile](../Makefile) to simplify building the project.

Use `make help` to list supported targets along with a short description:

```shell
$ make help

install-tools                            Install tools
clean-tools                              Remove installed tools
codegen-helm-docs                        Generate helm docs
codegen-schema-openapi                   Generate openapi schemas (v2 and v3)
codegen-all                              Generate all codegen
verify-schema-openapi                    Check openapi schemas are up to date
verify-helm-docs                         Check Helm charts are up to date
verify-codegen                           Verify all generated code and docs are up to date
build-clean                              Clean built files
build-frontend                           Build frontend
build-backend-assets                     Build backend assets
build-backend                            Build backend
build-all                                Build frontend and backend
ko-build                                 Build playground image (with ko)
ko-publish                               Build and publish playground image (with ko)
docker-build                             Build playground image (with docker)
run                                      Run locally
kind-create-cluster                      Create kind cluster
kind-delete-cluster                      Delete kind cluster
kind-load                                Build playground image and load it in kind cluster
kind-install                             Build image, load it in kind cluster and deploy playground helm chart
help                                     Shows the available commands
```

## Docker images

We use [ko](https://ko.build) to build docker images.

Run make `ko-build` to build local docker images.

You can build and publish docker images too by running:

```shell
REGISTRY=your_registry REPO=your_repo make ko-publish
```

## Building the frontend

Follow instruction in [frontend README](../frontend/README.md).

## Building the backend

Follow instruction in [backend README](../backend/README.md).
