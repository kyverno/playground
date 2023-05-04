export const PolicyTemplate = `apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: test-policy
spec:
  validationFailureAction: Audit  
  rules:
  - name: test-rule
    match: {}
    validate:
      message: ""`

export const ContextTemplate = `{
    "username": "",
    "groups": [],
    "role": [],
    "clusterrole": [],
    "operation": "CREATE"
}`

export const ResourceTemplate = `apiVersion: v1
kind: Pod
metadata:
  labels:
    run: nginx
  name: nginx
  namespace: default
spec:
  containers:
  - image: nginx
    name: nginx
    resources: {}`