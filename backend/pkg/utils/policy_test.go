package utils_test

import (
	"testing"

	"github.com/kyverno/playground/backend/pkg/resource/loader"
	"github.com/kyverno/playground/backend/pkg/utils"
)

const (
	singleResource string = `apiVersion: v1
kind: Namespace
metadata:
	name: prod-bus-app1
	labels:
	purpose: production`

	singlePolicy string = `apiVersion: kyverno.io/v1
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

	multiplePolicy string = `
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

	policyWithComment string = `
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
	loader, err := loader.New(nil, "1.27")
	if err != nil {
		t.Fatal(err)
	}
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
			if res, err := utils.LoadPolicies(loader, []byte(tt.policies)); (err != nil) != tt.wantErr {
				t.Errorf("loader.Policies() error = %v, wantErr %v", err, tt.wantErr)
			} else if len(res) != tt.wantLoaded {
				t.Errorf("loader.Policies() loaded amount = %v, wantLoaded %v", len(res), tt.wantLoaded)
			}
		})
	}
}
