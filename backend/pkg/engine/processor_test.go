package engine_test

import (
	"context"
	"testing"

	"github.com/kyverno/playground/backend/pkg/api"
	"github.com/kyverno/playground/backend/pkg/engine"
)

func Test_Processor(t *testing.T) {
	type fields struct {
		Policy    string
		Resources string
		Context   string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{
				Policy: `
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
              purpose: production`,
				Resources: `
apiVersion: v1
kind: Namespace
metadata:
  name: prod-bus-app1
  labels:
    purpose: production`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := api.Request{
				Policies:  tt.fields.Policy,
				Resources: tt.fields.Resources,
				Context:   tt.fields.Context,
			}

			params, err := engine.ParseParameters(r.Context)
			fatalOnError(t, "engine.ParseParameters", err)

			loader, err := api.NewLoader(params.Kubernetes.Version)
			fatalOnError(t, "api.NewLoader", err)

			resources, err := loader.Resources(r.Resources)
			fatalOnError(t, "loader.Resources", err)

			policies, err := loader.Policies(r.Policies)
			fatalOnError(t, "loader.Policies", err)

			processor, err := engine.NewProcessor(params)
			fatalOnError(t, "engine.NewProcessor", err)

			if _, err := processor.Run(context.TODO(), policies, resources); (err != nil) != tt.wantErr {
				t.Errorf("engineRequest.process() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func fatalOnError(t *testing.T, call string, err error) {
	if err == nil {
		return
	}

	t.Fatalf("%s error = %v", call, err)
}
