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

export const ConfigTemplate = `apiVersion: v1
kind: ConfigMap
metadata:
  name: kyverno
  namespace: kyverno
data:
  defaultRegistry: docker.io
  resourceFilters: |
    [*,kyverno,*]
    [Event,*,*]
    [*,kube-system,*]
    [*,kube-public,*]
    [*,kube-node-lease,*]
    [Node,*,*]
    [APIService,*,*]
    [TokenReview,*,*]
    [SubjectAccessReview,*,*]
    [SelfSubjectAccessReview,*,*]
    [Binding,*,*]
    [ReplicaSet,*,*]
    [AdmissionReport,*,*]
    [ClusterAdmissionReport,*,*]
    [BackgroundScanReport,*,*]
    [ClusterBackgroundScanReport,*,*]
    [ClusterRole,*,kyverno:*]
    [ClusterRoleBinding,*,kyverno:*]
    [ServiceAccount,kyverno,kyverno]
    [ConfigMap,kyverno,kyverno]
    [ConfigMap,kyverno,kyverno-metrics]
    [Deployment,kyverno,kyverno]
    [Job,kyverno,kyverno-hook-pre-delete]
    [NetworkPolicy,kyverno,kyverno]
    [PodDisruptionBudget,kyverno,kyverno]
    [Role,kyverno,kyverno:*]
    [RoleBinding,kyverno,kyverno:*]
    [Secret,kyverno,kyverno-svc.kyverno.svc.*]
    [Service,kyverno,kyverno-svc]
    [Service,kyverno,kyverno-svc-metrics]
    [ServiceMonitor,kyverno,kyverno-svc-service-monitor]
    [Pod,kyverno,kyverno-test]`

export const CustomResourceDefinitionsTemplate = `# apiVersion: apiextensions.k8s.io/v1
# kind: CustomResourceDefinition
# metadata:
#   annotations:
#     controller-gen.kubebuilder.io/version: v0.12.0
#   name: admirales.crew.testproject.org
# spec:
#   group: crew.testproject.org
#   names:
#     kind: Admiral
#     listKind: AdmiralList
#     plural: admirales
#     singular: admiral
#   scope: Cluster
#   versions:
#   - name: v1
#     schema:
#       openAPIV3Schema:
#         description: Admiral is the Schema for the admirales API
#         properties:
#           apiVersion:
#             description: 'APIVersion defines the versioned schema of this representation
#               of an object. Servers should convert recognized schemas to the latest
#               internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
#             type: string
#           kind:
#             description: 'Kind is a string value representing the REST resource this
#               object represents. Servers may infer this from the endpoint the client
#               submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
#             type: string
#           metadata:
#             type: object
#           spec:
#             description: AdmiralSpec defines the desired state of Admiral
#             properties:
#               foo:
#                 description: Foo is an example field of Admiral. Edit admiral_types.go
#                   to remove/update
#                 type: string
#             type: object
#           status:
#             description: AdmiralStatus defines the observed state of Admiral
#             type: object
#         type: object
#     served: true
#     storage: true
#     subresources:
#       status: {}`
