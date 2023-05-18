package api_test

import (
	"testing"

	"github.com/kyverno/playground/backend/pkg/api"
)

var loader, _ = api.NewLoader("1.27")

const (
	singleResource string = `apiVersion: v1
kind: Namespace
metadata:
  name: prod-bus-app1
  labels:
    purpose: production`

	multipleResources string = `
apiVersion: v1
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
      resources: {}
---
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: redis
  name: redis
  namespace: default
spec:
  containers:
    - image: redis
      name: redis
      resources: {}`

	resourceWithComment string = `
### POD ###
---
apiVersion: v1
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
)

func Test_LoadResources(t *testing.T) {
	tests := []struct {
		name       string
		resources  string
		wantLoaded int
		wantErr    bool
	}{
		{
			name:       "load no resource with empy string",
			resources:  "",
			wantLoaded: 0,
			wantErr:    false,
		},
		{
			name:       "load single resource",
			resources:  singleResource,
			wantLoaded: 1,
			wantErr:    false,
		},
		{
			name:       "load multiple resources",
			resources:  multipleResources,
			wantLoaded: 2,
			wantErr:    false,
		},
		{
			name:       "load resource with comment",
			resources:  resourceWithComment,
			wantLoaded: 1,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res, err := loader.Resources(tt.resources); (err != nil) != tt.wantErr {
				t.Errorf("loader.Resources() error = %v, wantErr %v", err, tt.wantErr)
			} else if len(res) != tt.wantLoaded {
				t.Errorf("loader.Resources() loaded amount = %v, wantLoaded %v", len(res), tt.wantLoaded)
			}
		})
	}
}

const (
	singlePolicy = `apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-ns-purpose-label
spec:
  validationFailureAction: Enforce
  rules:
  - name: require-ns-purpose-label
    match:
      any:
      - resources:
          kinds:
          - Namespace
    validate:
      message: "You must have label 'purpose' with value 'production' set on all new namespaces."
      pattern:
        metadata:
          labels:
            purpose: production`

	multiplePolicy = `
apiVersion: kyverno.io/v1
kind: Policy
metadata:
  name: require-ns-purpose-label
spec:
  validationFailureAction: Enforce
  rules:
  - name: require-ns-purpose-label
    match:
      any:
      - resources:
          kinds:
          - Namespace
    validate:
      message: "You must have label 'purpose' with value 'production' set on all new namespaces."
      pattern:
        metadata:
           labels:
           purpose: production
---
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-ns-purpose-label
spec:
  validationFailureAction: Enforce
  rules:
  - name: require-ns-purpose-label
    match:
      any:
      - resources:
          kinds:
          - Namespace
    validate:
      message: "You must have label 'purpose' with value 'production' set on all new namespaces."
      pattern:
        metadata:
          labels:
            purpose: production`

	policyWithComment = `
### Policy ###
---
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-ns-purpose-label
spec:
    validationFailureAction: Enforce
    rules:
    - name: require-ns-purpose-label
      match:
        any:
        - resources:
            kinds:
            - Namespace
      validate:
        message: "You must have label 'purpose' with value 'production' set on all new namespaces."
        pattern:
          metadata:
            labels:
              purpose: production`
)

func Test_LoadPolicies(t *testing.T) {
	tests := []struct {
		name       string
		policies   string
		wantLoaded int
		wantErr    bool
	}{
		{
			name:       "load no policy with empy string",
			policies:   "",
			wantLoaded: 0,
			wantErr:    false,
		},
		{
			name:       "load invalid string",
			policies:   singleResource,
			wantLoaded: 0,
			wantErr:    true,
		},
		{
			name:       "load single policy",
			policies:   singlePolicy,
			wantLoaded: 1,
			wantErr:    false,
		},
		{
			name:       "load multiple resources",
			policies:   multiplePolicy,
			wantLoaded: 2,
			wantErr:    false,
		},
		{
			name:       "load policy with comment",
			policies:   policyWithComment,
			wantLoaded: 1,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res, err := loader.Policies(tt.policies); (err != nil) != tt.wantErr {
				t.Errorf("loader.Policies() error = %v, wantErr %v", err, tt.wantErr)
			} else if len(res) != tt.wantLoaded {
				t.Errorf("loader.Policies() loaded amount = %v, wantLoaded %v", len(res), tt.wantLoaded)
			}
		})
	}
}
