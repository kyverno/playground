# Image tagging docs

This docs contains informations related to our image tagging strategy

## Dev images

All dev images published are tagged with their corresponding git commit hash.

Dev images are all images that are published outside of a release (including commits on `main` branch).

## Release images

Release images are tagged with their corresponding release tag.

Our release workflow is running `golreleaser` (which in turn runs `ko`) and is triggered when a tag is pushed to the repository.
The tag is expected to be in the form of a semantic release version (for example `v0.3.2-beta.1`).

## Latest tag

The `latest` tag is added only to stable releases, dev images or prerelease images never receive the `latest` tag.

Of course it is highly recommended to NOT use the `latest` tag, chances are that container flags will not be compatible from one version to another, it makes rollbacks difficult, etc...