export const PolicyTemplate = `apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-labels
spec:
  validationFailureAction: Audit
  rules:
    - name: check-for-labels
      match:
        any:
          - resources:
              kinds:
                - Pod
      validate:
        message: "label 'app.kubernetes.io/name' is required"
        pattern:
          metadata:
            labels:
              app.kubernetes.io/name: "?*"`

export const ContextTemplate = `kubernetes:
  version: '1.27'

context:
  username: ''
  groups: []
  roles: []
  clusterRoles: []
  namespaceLabels: {}
  operation: CREATE

variables: {}`

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
