package main

import "testing"

func Test_engineRequest_process(t *testing.T) {
	type fields struct {
		Policy   string
		Resource string
		Context  string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{
				Policy:   "",
				Resource: "",
			},
			wantErr: true,
		},
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
				Resource: `
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
			r := engineRequest{
				Policy:   tt.fields.Policy,
				Resource: tt.fields.Resource,
				Context:  tt.fields.Context,
			}
			if err := r.process(); (err != nil) != tt.wantErr {
				t.Errorf("engineRequest.process() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
