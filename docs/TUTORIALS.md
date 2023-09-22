# Tutorials docs

This docs contains informations related to the tutorials system used in the Playground

## tutorials.json

All tutorials are referenced in the [tutorials.json](../frontend/public/tutorials/tutorials.json) file.

This file contains an array of objects and defined the hierarchy and of the tutorials menu.

All properties can be inherited from the parent/overriden in the children.

| Property | Description |
|---|---|
| name                 | Menu item name                                                                  |
| url                  | URL where referenced files are looked up (can be an internet URL or local path) |
| color                | Menu item color                                                                 |
| subgroups            | Grouped children menu items                                                     |
| policies             | Children menu items                                                             |
| path                 | File path                                                                       |
| title                | Menu item title                                                                 |
| contextPath          | Path to context file                                                            |
| oldResourceFile      | Path to old resource file                                                       |
| clusterResourcesFile | Path to cluster resources file                                                  |
| exceptionsFile       | Path to exception file                                                          |

Example:

```json
    {
        "name": "Other",
        "url": "https://raw.githubusercontent.com/kyverno/policies/main/other",
        "policies": [
            {
                "path": "a/add-certificates-volume",
                "title": "Add Certificates as a Volume"
            },
            {
                "path": "a/add-default-resources",
                "title": "Add Default Resources"
            },
            {
                "path": "a/add-labels",
                "title": "Add Labels"
            },
            {
                "path": "a/allowed-annotations",
                "title": "Allowed Annotations"
            },
            {
                "path": "b-d/check-env-vars",
                "title": "Check Environment Variables"
            },
            {
                "path": "rec-req/require-base-image",
                "title": "Check Image Base"
            }
        ]
    }
```

Local files should be stored in the [tutorials folder](../frontend/public/tutorials/).
