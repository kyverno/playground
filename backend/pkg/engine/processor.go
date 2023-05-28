package engine

import (
	"context"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/go-logr/logr"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov1beta1 "github.com/kyverno/kyverno/api/kyverno/v1beta1"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/utils/store"
	"github.com/kyverno/kyverno/pkg/background/generate"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	"github.com/kyverno/kyverno/pkg/config"
	kyvernoengine "github.com/kyverno/kyverno/pkg/engine"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/kyverno/kyverno/pkg/engine/jmespath"
	"github.com/kyverno/kyverno/pkg/engine/mutate/patch"
	"github.com/kyverno/kyverno/pkg/engine/policycontext"
	"github.com/kyverno/kyverno/pkg/registryclient"
	jsonutils "github.com/kyverno/kyverno/pkg/utils/json"
	authenticationv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Processor struct {
	params        *Parameters
	engine        engineapi.Engine
	genController *generate.GenerateController
	config        config.Configuration
	jmesPath      jmespath.Interface
	cluster       bool
}

func (p *Processor) Run(
	ctx context.Context,
	policies []kyvernov1.PolicyInterface,
	resources []unstructured.Unstructured,
	oldResources []unstructured.Unstructured,
) (*Results, error) {
	if !p.cluster {
		if err := validateParams(p.params, policies); err != nil {
			return nil, err
		}
	}

	response := &Results{}

	for i, resource := range resources {
		oldResource := unstructured.Unstructured{}
		// @TODO does this make sense or is needed?
		if p.params.Context.Operation == kyvernov1.Delete {
			oldResource = resource
		}

		// @TODO should we also check for NS / kind / name or enforce the same order as resources?
		if len(oldResources) > i && p.params.Context.Operation == kyvernov1.Update {
			oldResource = oldResources[i]
		}

		// mutate
		for _, policy := range policies {
			result, res, err := p.mutate(ctx, policy, resource, oldResource)
			if err != nil {
				return nil, err
			}

			resource = res
			response.Mutation = append(response.Mutation, result)
		}

		// verify images
		for _, policy := range policies {
			result, res, err := p.verifyImages(ctx, policy, resource)
			if err != nil {
				return nil, err
			}

			resource = res
			response.ImageVerification = append(response.ImageVerification, result)
		}

		// validate
		for _, policy := range policies {
			result, err := p.validate(ctx, policy, resource, oldResource)
			if err != nil {
				return nil, err
			}

			response.Validation = append(response.Validation, result)
		}

		// generation
		for _, policy := range policies {
			result, err := p.generate(ctx, policy, resource, oldResource)
			if err != nil {
				return nil, err
			}

			response.Generation = append(response.Generation, result)
		}
	}

	return response, nil
}

func (p *Processor) mutate(ctx context.Context, policy kyvernov1.PolicyInterface, resource unstructured.Unstructured, oldResource unstructured.Unstructured) (Response, unstructured.Unstructured, error) {
	policyContext, err := p.newPolicyContext(policy, resource)
	if err != nil {
		return Response{}, resource, err
	}

	_ = policyContext.JSONContext().AddOldResource(oldResource.Object)

	response := p.engine.Mutate(ctx, policyContext.WithOldResource(oldResource))

	return ConvertResponse(response), response.PatchedResource, nil
}

func (p *Processor) verifyImages(ctx context.Context, policy kyvernov1.PolicyInterface, resource unstructured.Unstructured) (Response, unstructured.Unstructured, error) {
	policyContext, err := p.newPolicyContext(policy, resource)
	if err != nil {
		return Response{}, resource, err
	}

	response, verifiedImageData := p.engine.VerifyAndPatchImages(ctx, policyContext)
	// TODO: we apply patches manually because the engine doesn't
	patches := response.GetPatches()
	if !verifiedImageData.IsEmpty() {
		annotationPatches, err := verifiedImageData.Patches(len(resource.GetAnnotations()) != 0, logr.Discard())
		if err != nil {
			return Response{}, resource, err
		}

		// add annotation patches first
		patches = append(annotationPatches, patches...)
	}

	if len(patches) != 0 {
		patch := jsonutils.JoinPatches(patch.ConvertPatches(patches...)...)
		decoded, err := jsonpatch.DecodePatch(patch)
		if err != nil {
			return Response{}, resource, err
		}
		options := &jsonpatch.ApplyOptions{SupportNegativeIndices: true, AllowMissingPathOnRemove: true, EnsurePathExistsOnAdd: true}
		resourceBytes, err := resource.MarshalJSON()
		if err != nil {
			return Response{}, resource, err
		}
		patchedResourceBytes, err := decoded.ApplyWithOptions(resourceBytes, options)
		if err != nil {
			return Response{}, resource, err
		}
		if err := response.PatchedResource.UnmarshalJSON(patchedResourceBytes); err != nil {
			return Response{}, resource, err
		}
	}

	return ConvertResponse(response), response.PatchedResource, nil
}

func (p *Processor) validate(ctx context.Context, policy kyvernov1.PolicyInterface, resource unstructured.Unstructured, oldResource unstructured.Unstructured) (Response, error) {
	policyContext, err := p.newPolicyContext(policy, resource)
	if err != nil {
		return Response{}, err
	}

	_ = policyContext.JSONContext().AddOldResource(oldResource.Object)

	response := p.engine.Validate(ctx, policyContext.WithOldResource(oldResource))

	return ConvertResponse(response), nil
}

func (p *Processor) generate(ctx context.Context, policy kyvernov1.PolicyInterface, resource unstructured.Unstructured, oldResource unstructured.Unstructured) (Response, error) {
	policyContext, err := p.newPolicyContext(policy, resource)
	if err != nil {
		return Response{}, err
	}

	_ = policyContext.JSONContext().AddOldResource(oldResource.Object)

	policyContext = policyContext.WithOldResource(oldResource)

	response := p.engine.Generate(ctx, policyContext)
	if len(response.PolicyResponse.Rules) == 0 {
		return ConvertResponse(response), nil
	}

	gr := toGenerateRequest(policy, resource)

	var newRuleResponse []engineapi.RuleResponse
	for _, rule := range response.PolicyResponse.Rules {
		genRes, err := p.genController.ApplyGeneratePolicy(logr.Discard(), policyContext, gr, []string{rule.Name()})
		if err != nil {
			return Response{}, err
		}
		if len(genRes) == 0 {
			continue
		}
		unstrGenResource, err := p.genController.GetUnstrResource(genRes[0])
		if err != nil {
			return Response{}, err
		}

		// cleanup metadata
		if meta, ok := unstrGenResource.Object["metadata"]; ok {
			delete(meta.(map[string]any), "managedFields")
		}

		newRuleResponse = append(newRuleResponse, *rule.WithGeneratedResource(*unstrGenResource))
	}
	response.PolicyResponse.Rules = newRuleResponse

	return ConvertResponse(response), nil
}

func (p *Processor) newPolicyContext(policy kyvernov1.PolicyInterface, resource unstructured.Unstructured) (*policycontext.PolicyContext, error) {
	context, err := policycontext.NewPolicyContext(
		p.jmesPath,
		resource,
		p.params.Context.Operation,
		&kyvernov1beta1.RequestInfo{
			AdmissionUserInfo: authenticationv1.UserInfo{
				Username: p.params.Context.Username,
				Groups:   p.params.Context.Groups,
			},
			Roles:        p.params.Context.Roles,
			ClusterRoles: p.params.Context.ClusterRoles,
		},
		p.config,
	)
	if err != nil {
		return nil, err
	}

	for k, v := range p.params.Variables {
		err = context.JSONContext().AddVariable(k, v)
		if err != nil {
			return nil, err
		}
	}

	context = context.
		WithPolicy(policy).
		WithNamespaceLabels(p.params.Context.NamespaceLabels)

	return context, nil
}

func validateParams(params *Parameters, policies []kyvernov1.PolicyInterface) error {
	if params == nil {
		return nil
	}

	for _, policy := range policies {
		for _, rule := range policy.GetSpec().Rules {
			for _, variable := range rule.Context {
				if variable.APICall == nil && variable.ConfigMap == nil {
					continue
				}

				if _, ok := params.Variables[variable.Name]; !ok {
					return fmt.Errorf("Variable %s is not defined in the context", variable.Name)
				}
			}
		}
	}

	return nil
}

func toGenerateRequest(policy kyvernov1.PolicyInterface, resource unstructured.Unstructured) kyvernov1beta1.UpdateRequest {
	return kyvernov1beta1.UpdateRequest{
		Spec: kyvernov1beta1.UpdateRequestSpec{
			Type:   kyvernov1beta1.Generate,
			Policy: policy.GetName(),
			Resource: kyvernov1.ResourceSpec{
				Kind:       resource.GetKind(),
				Namespace:  resource.GetNamespace(),
				Name:       resource.GetName(),
				APIVersion: resource.GetAPIVersion(),
			},
		},
	}
}

func newEngine(
	cfg config.Configuration,
	jp jmespath.Interface,
	client dclient.Interface,
	cmResolver engineapi.ConfigmapResolver,
	exceptionSelector engineapi.PolicyExceptionSelector,
) (engineapi.Engine, error) {
	rclient, err := registryclient.New(registryclient.WithLocalKeychain())
	if err != nil {
		return nil, err
	}

	store.SetMock(true)
	store.SetRegistryAccess(true)

	factory := store.ContextLoaderFactory(nil)
	if cmResolver != nil {
		factory = ContextLoaderFactory(cmResolver)
	}

	return kyvernoengine.NewEngine(
		cfg,
		config.NewDefaultMetricsConfiguration(),
		jp,
		client,
		rclient,
		factory,
		exceptionSelector,
	), nil
}

func NewProcessor(
	params *Parameters,
	kyvernoConfig *corev1.ConfigMap,
	dClient dclient.Interface,
	cmResolver engineapi.ConfigmapResolver,
	exceptionSelector engineapi.PolicyExceptionSelector,
) (*Processor, error) {
	cfg := config.NewDefaultConfiguration(false)
	if kyvernoConfig != nil {
		cfg.Load(kyvernoConfig)
	}

	jp := jmespath.New(cfg)
	cluster := false

	engine, err := newEngine(cfg, jp, dClient, cmResolver, exceptionSelector)
	if err != nil {
		return nil, err
	}

	if dClient == nil {
		dClient = dclient.NewEmptyFakeClient()
	} else {
		cluster = true
	}

	contr := generate.NewGenerateController(dClient, nil, nil, engine, nil, nil, nil, nil, cfg, nil, logr.Discard(), jp)

	return &Processor{
		params:        params,
		engine:        engine,
		genController: contr,
		config:        cfg,
		jmesPath:      jp,
		cluster:       cluster,
	}, nil
}

func ContextLoaderFactory(cmResolver engineapi.ConfigmapResolver) engineapi.ContextLoaderFactory {
	return func(policy kyvernov1.PolicyInterface, rule kyvernov1.Rule) engineapi.ContextLoader {
		inner := engineapi.DefaultContextLoaderFactory(cmResolver)

		return inner(policy, rule)
	}
}
