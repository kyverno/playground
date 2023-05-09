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

export const ContextTemplate = `username: ''
groups: []
roles: []
clusterRoles: []
namespaceLabels: {}
operation: 'CREATE'`

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

export const PodTemplate = `apiVersion: v1
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

export const DeploymentTemplate = `apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx
        resources: {}`

export const ServiceTemplate = `apiVersion: v1
kind: Service
metadata:
  labels:
    app: nginx
  name: nginx
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    run: nginx
  type: ClusterIP
`

export const IngressTemplate = `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: nginx
  namespace: default
spec:
  ingressClassName: nginx
  rules:
  - host: app.example.com
    http:
      paths:
      - backend:
          service:
            name: nginx
            port:
              number: 80
        path: /(.*)
        pathType: ImplementationSpecific
  tls:
  - hosts:
    - app.example.com
    secretName: app-tls
`

export const ConfigMapTemplate = `apiVersion: v1
data:
  key: label
kind: ConfigMap
metadata:
  name: app-config
`

export const SecretTemplate = `apiVersion: v1
data:
  key: bGFiZWw=
kind: Secret
metadata:
  name: app-secret
`
