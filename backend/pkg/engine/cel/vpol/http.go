package vpol

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno-authz/pkg/cel/libs/authz/http"
	"github.com/kyverno/kyverno-authz/pkg/engine"
	vpolcompiler "github.com/kyverno/kyverno-authz/pkg/engine/compiler"
	"github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/kyverno/sdk/core"
	"github.com/kyverno/sdk/core/dispatchers"
	"github.com/kyverno/sdk/core/handlers"
	"github.com/kyverno/sdk/core/resulters"
	"github.com/kyverno/sdk/extensions/policy"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"

	"github.com/kyverno/playground/backend/pkg/engine/models"
)

func HTTPProcess(ctx context.Context, resource *http.CheckRequest, vpols []*v1beta1.ValidatingPolicy) ([]models.Response, error) {
	compiler := vpolcompiler.NewCompiler[dynamic.Interface, *http.CheckRequest, *http.CheckResponse]()
	results := make([]models.Response, 0)

	for _, vpol := range vpols {
		eng := core.NewEngine(
			NewHTTPSource(vpol, compiler),
			handlers.Handler(
				dispatchers.Sequential(
					policy.EvaluatorFactory[engine.HTTPPolicy](),
					func(ctx context.Context, fc core.FactoryContext[engine.HTTPPolicy, dynamic.Interface, *http.CheckRequest]) core.Breaker[engine.HTTPPolicy, *http.CheckRequest, policy.Evaluation[*http.CheckResponse]] {
						return core.MakeBreakerFunc(func(_ context.Context, _ engine.HTTPPolicy, _ *http.CheckRequest, out policy.Evaluation[*http.CheckResponse]) bool {
							return out.Result != nil
						})
					},
				),
				func(ctx context.Context, fc core.FactoryContext[engine.HTTPPolicy, dynamic.Interface, *http.CheckRequest]) core.Resulter[engine.HTTPPolicy, *http.CheckRequest, policy.Evaluation[*http.CheckResponse], policy.Evaluation[*http.CheckResponse]] {
					return resulters.NewFirst[engine.HTTPPolicy, *http.CheckRequest](func(out policy.Evaluation[*http.CheckResponse]) bool {
						return out.Result != nil || out.Error != nil
					})
				},
			),
		)
		evaluation := eng.Handle(ctx, nil, resource)

		var status api.RuleStatus
		var message string
		resource := unstructured.Unstructured{Object: make(map[string]any)}

		if evaluation.Result == nil && evaluation.Error == nil {
			message = "request does not match"
			status = api.RuleStatusSkip
		}

		if evaluation.Result != nil {
			var content []byte
			var err error

			if ok := evaluation.Result.Ok; ok != nil {
				content, err = json.Marshal(evaluation.Result.Ok)
				if err != nil {
					continue
				}
				status = api.RuleStatusPass
				message = "request allowed"
			} else if denied := evaluation.Result.Denied; denied != nil {
				content, err = json.Marshal(denied)
				if err != nil {
					continue
				}
				status = api.RuleStatusFail
				message = fmt.Sprintf("request denied, reason: %s", denied.Reason)
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

func NewHTTPSource(vpol v1beta1.ValidatingPolicyLike, compiler engine.Compiler[engine.HTTPPolicy]) core.Source[policy.Policy[dynamic.Interface, *http.CheckRequest, *http.CheckResponse]] {
	if v, ok := vpol.(*v1beta1.ValidatingPolicy); ok {
		p, err := compiler.Compile(v)
		if err != nil {
			return nil
		}
		return core.MakeSource(p)
	}

	return nil
}
