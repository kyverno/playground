package engine_test

// import (
// 	"context"
// 	"testing"

// 	"github.com/kyverno/playground/backend/pkg/engine"
// 	"github.com/kyverno/playground/backend/pkg/resource/loader"
// 	"github.com/kyverno/playground/backend/pkg/utils"
// )

// func Test_Processor(t *testing.T) {
// 	type fields struct {
// 		Policy    string
// 		Resources string
// 		Context   string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr bool
// 	}{
// 		{
// 			fields: fields{
// 				Policy: `
// apiVersion: kyverno.io/v1
// kind: ClusterPolicy
// metadata:
//   name: require-ns-purpose-label
// spec:
//   validationFailureAction: Enforce
//   rules:
//     - name: require-ns-purpose-label
//       match:
//         any:
//           - resources:
//               kinds:
//                 - Namespace
//       validate:
//         message: "You must have label 'purpose' with value 'production' set on all new namespaces."
//         pattern:
//           metadata:
//             labels:
//               purpose: production`,
// 				Resources: `
// apiVersion: v1
// kind: Namespace
// metadata:
//   name: prod-bus-app1
//   labels:
//     purpose: production`,
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := api.EngineRequest{
// 				Policies:  tt.fields.Policy,
// 				Resources: tt.fields.Resources,
// 				Context:   tt.fields.Context,
// 			}

// 			params, err := engine.ParseParameters(r.Context)
// 			fatalOnError(t, "engine.ParseParameters", err)

// 			l, err := loader.NewLocal(params.Kubernetes.Version)
// 			fatalOnError(t, "loader.New", err)

// 			resources, err := loader.LoadResources(l, []byte(r.Resources))
// 			fatalOnError(t, "loader.LoadResources", err)

// 			policies, err := utils.LoadPolicies(l, []byte(r.Policies))
// 			fatalOnError(t, "loader.LoadPolicies", err)

// 			processor, err := engine.NewProcessor(params, nil, nil, nil, nil)
// 			fatalOnError(t, "engine.NewProcessor", err)

// 			if _, err := processor.Run(context.TODO(), policies, resources); (err != nil) != tt.wantErr {
// 				t.Errorf("engineRequest.process() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func fatalOnError(t *testing.T, call string, err error) {
// 	if err == nil {
// 		return
// 	}

// 	t.Fatalf("%s error = %v", call, err)
// }
