# Release docs

This docs contains informations related to releasing the Playground

## Create a release

Creating a release can be done by pushing a tag to the GitHub repository (begining with `v`).

The [release workflow](../.github/workflows/release.yaml) will take care of creating the GitHub release and will publish docker images.

```shell
VERSION="v0.1.0"
TAG=$VERSION

git tag $TAG -m "tag $TAG" -a
git push origin $TAG
```

## Create an Helm release

Creating an Helm release can be done by pushing a tag to the GitHub repository (begining with `kyverno-playground-chart-v`).

The [helm workflow](../.github/workflows/helm.yaml) will take care of creating the Helm release and will publish it to https://kyverno.github.io/playground.

```shell
VERSION="v0.1.0"
TAG=kyverno-playground-chart-$VERSION

git tag $TAG -m "tag $TAG" -a
git push origin $TAG
```
