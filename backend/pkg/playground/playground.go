package playground

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	yamlutils "github.com/kyverno/kyverno/ext/yaml"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"google.golang.org/protobuf/encoding/protojson"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/crd"
	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/resource"
	"github.com/kyverno/playground/backend/pkg/serve"
)

// The Playground has no policy evaluation of its own: every request is evaluated by a
// `kyverno serve --allow-request-policies` instance through its v1alpha1 HTTP API, so the
// Playground, `kyverno apply` and `kyverno serve` share one evaluation path.

var (
	clientMu sync.RWMutex
	client   = serve.NewClient(serve.URLFromEnv())
)

// SetServeURL points the Playground at the kyverno serve instance that evaluates requests.
func SetServeURL(url string) {
	clientMu.Lock()
	defer clientMu.Unlock()
	client = serve.NewClient(url)
}

func serveClient() *serve.Client {
	clientMu.RLock()
	defer clientMu.RUnlock()
	return client
}

const (
	modeKubernetes = "Kubernetes"
	modeJSON       = "JSON"
	modeHTTP       = "HTTP"
	modeEnvoy      = "Envoy"
)

// Run evaluates request with kyverno serve. The cluster and CRD configuration parameters are kept
// for API compatibility; kyverno serve evaluates offline against the objects in the request.
func Run(ctx context.Context, _ cluster.Cluster, request *EngineRequest, _ crd.APIConfiguration) (*EngineResponse, error) {
	params, err := request.LoadParameters()
	if err != nil {
		return nil, err
	}
	if err := checkSupported(params); err != nil {
		return nil, err
	}
	policies, err := parsePolicies(request.Policies)
	if err != nil {
		return nil, err
	}
	mode := policies.mode()
	if mode == "" || strings.TrimSpace(request.Resources) == "" {
		return &EngineResponse{}, nil
	}
	payloads, resources, err := splitPayloads(mode, request.Resources)
	if err != nil {
		return nil, fmt.Errorf("unable to load resources: %w", err)
	}

	base := serve.EvaluationRequest{
		Mode:       mode,
		Policies:   joinDocuments(request.Policies, request.VAPBindings),
		Exceptions: request.PolicyExceptions,
	}
	var clusterObjects []unstructured.Unstructured
	if mode == modeKubernetes {
		base.Operation = string(params.Context.Operation)
		base.CustomResourceDefinitions = request.CustomResourceDefinitions
		base.OldPayload = request.OldResources
		// Playground cluster resources are the existing cluster: namespace lookups, CEL
		// resource.get/resource.list (through a CLI Context) and mutateExisting targets read them.
		existing := request.ClusterResources
		base.ClusterResources = existing
		base.TargetResources = existing
		if base.Context, err = buildContext(existing, params.ImageData, request.ImageData); err != nil {
			return nil, err
		}
		objects, err := yamlObjects(existing)
		if err != nil {
			return nil, fmt.Errorf("unable to load cluster resources: %w", err)
		}
		for _, o := range objects {
			clusterObjects = append(clusterObjects, unstructured.Unstructured{Object: o})
		}
	}

	results := &models.Results{
		Mutation:          []models.Response{},
		ImageVerification: []models.Response{},
		Deletion:          []models.Response{},
		Validation:        []models.Response{},
		Generation:        []models.Response{},
	}
	for i, payload := range payloads {
		req := base
		req.Payload = payload
		resp, err := serveClient().Evaluate(ctx, req)
		if err != nil {
			var serveErr *serve.Error
			if errors.As(err, &serveErr) && serveErr.Code < 500 {
				return nil, errors.New(serveErr.Message)
			}
			return nil, err
		}
		var payloadObject *unstructured.Unstructured
		if mode != modeKubernetes && i < len(resources) {
			payloadObject = &resources[i]
		}
		addResults(results, resp, policies, append(cloneObjects(resources), clusterObjects...), len(resources), payloadObject)
	}
	out := &EngineResponse{Results: results}
	if mode == modeKubernetes || mode == modeJSON {
		out.Resources = resources
	}
	return out, nil
}

func cloneObjects(in []unstructured.Unstructured) []unstructured.Unstructured {
	return append([]unstructured.Unstructured(nil), in...)
}

// checkSupported rejects request settings kyverno serve v1alpha1 cannot honor, instead of silently
// evaluating without them.
func checkSupported(params *models.Parameters) error {
	var unsupported []string
	if len(params.Variables) > 0 {
		unsupported = append(unsupported, "context variables")
	}
	c := params.Context
	if c.Username != "" || len(c.Groups) > 0 || len(c.Roles) > 0 || len(c.ClusterRoles) > 0 {
		unsupported = append(unsupported, "user info (username, groups, roles, clusterRoles)")
	}
	if len(c.NamespaceLabels) > 0 {
		unsupported = append(unsupported, "context namespaceLabels (add the Namespace to cluster resources instead)")
	}
	if len(unsupported) > 0 {
		return fmt.Errorf("not supported by the kyverno serve v1alpha1 evaluation backend: %s", strings.Join(unsupported, ", "))
	}
	return nil
}

func joinDocuments(docs ...string) string {
	var parts []string
	for _, d := range docs {
		if strings.TrimSpace(d) != "" {
			parts = append(parts, strings.TrimSpace(d))
		}
	}
	return strings.Join(parts, "\n---\n")
}

// buildContext converts cluster objects and image data into a cli.kyverno.io/v1alpha1 Context.
func buildContext(clusterResources string, imageData ...map[string]models.ImageData) (string, error) {
	spec := map[string]any{}
	if strings.TrimSpace(clusterResources) != "" {
		objects, err := yamlObjects(clusterResources)
		if err != nil {
			return "", fmt.Errorf("unable to load cluster resources: %w", err)
		}
		if len(objects) > 0 {
			spec["resources"] = objects
		}
	}
	var images []map[string]any
	for _, m := range imageData {
		for _, img := range m {
			images = append(images, map[string]any{
				"image":         img.Image,
				"resolvedImage": img.ResolvedImage,
				"registry":      img.Registry,
				"repository":    img.Repository,
				"identifier":    img.Identifier,
				"manifest":      img.Manifest,
				"configData":    img.ConfigData,
			})
		}
	}
	if len(images) > 0 {
		spec["images"] = images
	}
	if len(spec) == 0 {
		return "", nil
	}
	b, err := yaml.Marshal(map[string]any{
		"apiVersion": "cli.kyverno.io/v1alpha1",
		"kind":       "Context",
		"metadata":   map[string]any{"name": "playground"},
		"spec":       spec,
	})
	return string(b), err
}

func yamlObjects(content string) ([]map[string]any, error) {
	docs, err := yamlutils.SplitDocuments([]byte(content))
	if err != nil {
		return nil, err
	}
	var out []map[string]any
	for _, doc := range docs {
		var obj map[string]any
		if err := yaml.Unmarshal(doc, &obj); err != nil {
			return nil, err
		}
		if len(obj) == 0 {
			continue
		}
		if items, ok := obj["items"].([]any); ok && strings.HasSuffix(fmt.Sprint(obj["kind"]), "List") {
			for _, item := range items {
				if m, ok := item.(map[string]any); ok {
					out = append(out, m)
				}
			}
			continue
		}
		out = append(out, obj)
	}
	return out, nil
}

// splitPayloads turns the Playground "resources" editor content into kyverno serve payloads: one
// request with every Kubernetes resource, or one request per JSON/HTTP/Envoy payload (a JSON
// array holds several). It also returns the decoded objects shown back in the UI.
func splitPayloads(mode, content string) ([]string, []unstructured.Unstructured, error) {
	switch mode {
	case modeKubernetes:
		objects, err := yamlObjects(content)
		if err != nil {
			return nil, nil, err
		}
		resources := make([]unstructured.Unstructured, 0, len(objects))
		for _, o := range objects {
			resources = append(resources, unstructured.Unstructured{Object: o})
		}
		return []string{content}, resources, nil
	case modeJSON:
		resources, err := resource.LoadJSON(content)
		if err != nil {
			return nil, nil, err
		}
		payloads := make([]string, 0, len(resources))
		for _, r := range resources {
			b, err := json.Marshal(r.Object)
			if err != nil {
				return nil, nil, err
			}
			payloads = append(payloads, string(b))
		}
		return payloads, resources, nil
	case modeEnvoy:
		requests, err := resource.LoadEnvyRequests(content)
		if err != nil {
			return nil, nil, err
		}
		payloads := make([]string, 0, len(requests))
		for _, r := range requests {
			b, err := protojson.Marshal(r)
			if err != nil {
				return nil, nil, err
			}
			payloads = append(payloads, string(b))
		}
		return payloads, nil, nil
	case modeHTTP:
		requests, err := resource.LoadHTTPRequests(content)
		if err != nil {
			return nil, nil, err
		}
		payloads := make([]string, 0, len(requests))
		for _, r := range requests {
			b, err := json.Marshal(r)
			if err != nil {
				return nil, nil, err
			}
			payloads = append(payloads, string(b))
		}
		return payloads, nil, nil
	}
	return nil, nil, fmt.Errorf("unsupported mode %s", mode)
}

// addResults groups kyverno serve rule results by policy, resource and rule type into the
// Playground's per-type response lists, ordered like the objects in the request.
func addResults(out *models.Results, resp *serve.EvaluationResponse, policies policyIndex, objects []unstructured.Unstructured, payloadCount int, payload *unstructured.Unstructured) {
	type key struct {
		policy   serve.ObjectRef
		resource serve.ObjectRef
		ruleType string
	}
	type group struct {
		key      key
		position int
		seq      int
		response *models.Response
	}
	var groups []*group
	byKey := map[key]*group{}
	for _, r := range resp.Results {
		var resRef serve.ObjectRef
		if r.Resource != nil {
			resRef = *r.Resource
		}
		k := key{policy: r.Policy, resource: resRef, ruleType: r.Type}
		g, ok := byKey[k]
		if !ok {
			g = &group{key: k, position: len(objects), seq: len(groups), response: &models.Response{Policy: policies.describe(r.Policy)}}
			if idx, obj := findResource(objects, r.Resource, payload); obj != nil {
				g.position = idx
				g.response.Resource = *obj
				if b, err := yaml.Marshal(obj.Object); err == nil {
					g.response.OriginalResource = string(b)
					g.response.PatchedResource = string(b)
				}
			}
			byKey[k] = g
			groups = append(groups, g)
		}
		patched := r.PatchedResource
		if r.PatchedTarget != nil {
			patched = r.PatchedTarget
		}
		if patched != nil {
			if b, err := yaml.Marshal(patched); err == nil {
				g.response.PatchedResource = string(b)
			}
		}
		g.response.PolicyResponse.Rules = append(g.response.PolicyResponse.Rules, convertRule(r))
	}
	// payload resources keep their order; results on other objects (mutateExisting targets) follow,
	// ordered by namespace and name
	sort.SliceStable(groups, func(i, j int) bool {
		pi, pj := groups[i].position, groups[j].position
		if pi >= payloadCount && pj >= payloadCount {
			ri, rj := groups[i].key.resource, groups[j].key.resource
			if ri.Namespace != rj.Namespace {
				return ri.Namespace < rj.Namespace
			}
			if ri.Name != rj.Name {
				return ri.Name < rj.Name
			}
			return groups[i].seq < groups[j].seq
		}
		if pi != pj {
			return pi < pj
		}
		return groups[i].seq < groups[j].seq
	})
	for _, g := range groups {
		switch engineapi.RuleType(g.key.ruleType) {
		case engineapi.Mutation:
			out.Mutation = append(out.Mutation, *g.response)
		case engineapi.Generation:
			out.Generation = append(out.Generation, *g.response)
		case engineapi.ImageVerify:
			out.ImageVerification = append(out.ImageVerification, *g.response)
		case engineapi.Deletion:
			out.Deletion = append(out.Deletion, *g.response)
		default:
			out.Validation = append(out.Validation, *g.response)
		}
	}
}

func convertRule(r serve.RuleResult) models.RuleResponse {
	name := r.Rule
	if name == "" {
		name = r.Policy.Name
	}
	rule := models.RuleResponse{
		Name:       name,
		RuleType:   engineapi.RuleType(r.Type),
		Message:    r.Message,
		Status:     engineapi.RuleStatus(r.Status),
		Properties: r.Properties,
	}
	if r.PatchedTarget != nil {
		rule.PatchedTarget = &unstructured.Unstructured{Object: r.PatchedTarget}
	}
	switch len(r.GeneratedResources) {
	case 0:
	case 1:
		if b, err := yaml.Marshal(r.GeneratedResources[0]); err == nil {
			rule.GeneratedResource = string(b)
		}
	default:
		if b, err := yaml.Marshal(r.GeneratedResources); err == nil {
			rule.GeneratedResource = string(b)
		}
	}
	return rule
}

func findResource(objects []unstructured.Unstructured, ref *serve.ObjectRef, payload *unstructured.Unstructured) (int, *unstructured.Unstructured) {
	if ref == nil {
		return 0, payload
	}
	for i := range objects {
		r := &objects[i]
		if r.GetAPIVersion() == ref.APIVersion && r.GetKind() == ref.Kind && r.GetNamespace() == ref.Namespace && r.GetName() == ref.Name {
			return i, r
		}
	}
	return len(objects), payload
}
