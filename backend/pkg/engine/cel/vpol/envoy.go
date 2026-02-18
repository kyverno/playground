package vpol

import (
	"context"
	"encoding/json"
	"fmt"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno-authz/pkg/engine"
	vpolcompiler "github.com/kyverno/kyverno-authz/pkg/engine/compiler"
	"github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/sdk/core"
	"github.com/kyverno/sdk/core/dispatchers"
	"github.com/kyverno/sdk/core/handlers"
	"github.com/kyverno/sdk/core/resulters"
	"github.com/kyverno/sdk/extensions/policy"
	"google.golang.org/protobuf/encoding/protojson"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

func EnvoyProcess(ctx context.Context, resource *authv3.CheckRequest, vpols []*v1beta1.ValidatingPolicy) ([]models.Response, error) {
	compiler := vpolcompiler.NewCompiler[dynamic.Interface, *authv3.CheckRequest, *authv3.CheckResponse]()
	results := make([]models.Response, 0)

	for _, vpol := range vpols {
		eng := core.NewEngine(
			NewSingleSource(vpol, compiler),
			handlers.Handler(
				dispatchers.Sequential(
					policy.EvaluatorFactory[engine.EnvoyPolicy](),
					func(ctx context.Context, fc core.FactoryContext[engine.EnvoyPolicy, dynamic.Interface, *authv3.CheckRequest]) core.Breaker[engine.EnvoyPolicy, *authv3.CheckRequest, policy.Evaluation[*authv3.CheckResponse]] {
						return core.MakeBreakerFunc(func(_ context.Context, _ engine.EnvoyPolicy, _ *authv3.CheckRequest, out policy.Evaluation[*authv3.CheckResponse]) bool {
							return out.Result != nil
						})
					},
				),
				func(ctx context.Context, fc core.FactoryContext[engine.EnvoyPolicy, dynamic.Interface, *authv3.CheckRequest]) core.Resulter[engine.EnvoyPolicy, *authv3.CheckRequest, policy.Evaluation[*authv3.CheckResponse], policy.Evaluation[*authv3.CheckResponse]] {
					return resulters.NewFirst[engine.EnvoyPolicy, *authv3.CheckRequest](func(out policy.Evaluation[*authv3.CheckResponse]) bool {
						return out.Result != nil || out.Error != nil
					})
				},
			),
		)

		evaluation := eng.Handle(ctx, nil, resource)

		if evaluation.Result == nil && evaluation.Error == nil {
			continue
		}

		var status api.RuleStatus
		var message string

		var resource unstructured.Unstructured
		if evaluation.Result != nil {
			var content []byte
			var err error

			if ok := evaluation.Result.GetOkResponse(); ok != nil {
				content, err = protojson.Marshal(evaluation.Result.GetOkResponse())
				if err != nil {
					continue
				}
				status = api.RuleStatusPass
				message = "request allowed"
			} else if denied := evaluation.Result.GetDeniedResponse(); denied != nil {
				content, err = protojson.Marshal(denied)
				if err != nil {
					continue
				}
				status = api.RuleStatusFail
				message = fmt.Sprintf("request denied with status code %d (\"%s\")", denied.Status.Code, denied.Status.Code.String())
			}

			var payload map[string]any
			if err := json.Unmarshal([]byte(content), &payload); err != nil {
				return nil, err
			}

			resource = unstructured.Unstructured{
				Object: payload,
			}
		}

		if evaluation.Error != nil {
			message = evaluation.Error.Error()
			status = api.RuleStatusError
		}

		results = append(results, models.Response{
			Policy: models.Policy{
				APIVersion:  vpol.APIVersion,
				Name:        vpol.Name,
				Namespace:   vpol.Namespace,
				Mode:        vpol.Spec.EvaluationMode(),
				Kind:        vpol.Kind,
				Labels:      vpol.Labels,
				Annotations: vpol.Annotations,
			},
			PolicyResponse: models.PolicyResponse{
				Rules: []models.RuleResponse{
					{
						Name:     vpol.Name,
						RuleType: api.Validation,
						Message:  message,
						Status:   status,
					},
				},
			},
			Resource: resource,
		})
	}

	return results, nil
}

func NewSource(vpols []v1beta1.ValidatingPolicyLike, compiler engine.Compiler[policy.Policy[dynamic.Interface, *authv3.CheckRequest, *authv3.CheckResponse]]) core.Source[policy.Policy[dynamic.Interface, *authv3.CheckRequest, *authv3.CheckResponse]] {
	polcies := make([]policy.Policy[dynamic.Interface, *authv3.CheckRequest, *authv3.CheckResponse], 0, len(vpols))

	for _, vpol := range vpols {
		if v, ok := vpol.(*v1beta1.ValidatingPolicy); ok {
			p, err := compiler.Compile(v)
			if err != nil {
				continue
			}

			polcies = append(polcies, p)
		}
	}

	return core.MakeSource(polcies...)
}

func NewSingleSource(vpol v1beta1.ValidatingPolicyLike, compiler engine.Compiler[policy.Policy[dynamic.Interface, *authv3.CheckRequest, *authv3.CheckResponse]]) core.Source[policy.Policy[dynamic.Interface, *authv3.CheckRequest, *authv3.CheckResponse]] {
	if v, ok := vpol.(*v1beta1.ValidatingPolicy); ok {
		p, err := compiler.Compile(v)
		if err != nil {
			return nil
		}
		return core.MakeSource(p)
	}

	return nil
}
