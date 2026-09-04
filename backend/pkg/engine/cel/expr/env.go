package expr

import (
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/ext"
	"github.com/kyverno/kyverno/pkg/cel/compiler"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	"github.com/kyverno/sdk/extensions/cel/libs/globalcontext"
	"github.com/kyverno/sdk/extensions/cel/libs/gzip"
	"github.com/kyverno/sdk/extensions/cel/libs/hash"
	"github.com/kyverno/sdk/extensions/cel/libs/http"
	"github.com/kyverno/sdk/extensions/cel/libs/image"
	"github.com/kyverno/sdk/extensions/cel/libs/imagedata"
	"github.com/kyverno/sdk/extensions/cel/libs/json"
	"github.com/kyverno/sdk/extensions/cel/libs/math"
	"github.com/kyverno/sdk/extensions/cel/libs/random"
	"github.com/kyverno/sdk/extensions/cel/libs/resource"
	"github.com/kyverno/sdk/extensions/cel/libs/time"
	"github.com/kyverno/sdk/extensions/cel/libs/transform"
	"github.com/kyverno/sdk/extensions/cel/libs/user"
	"github.com/kyverno/sdk/extensions/cel/libs/x509"
	"github.com/kyverno/sdk/extensions/cel/libs/yaml"
	apiservercel "k8s.io/apiserver/pkg/cel"
	"k8s.io/apiserver/pkg/cel/environment"
)

// newEnv builds a CEL environment for evaluating standalone expressions with
// the same variables (object, oldObject, request, namespaceObject, variables)
// and Kyverno CEL libraries (resource fetching, HTTP calls, global context,
// image data, etc.) that a ValidatingPolicy would have available. It mirrors
// the (unexported) environment construction in kyverno's
// pkg/cel/policies/vpol/compiler package, since that logic isn't reusable
// across module boundaries.
func newEnv(libsctx libs.Context, namespace string) (*cel.Env, error) {
	baseOpts := compiler.DefaultEnvOptionsWithCompat()
	baseOpts = append(baseOpts,
		cel.Variable(compiler.NamespaceObjectKey, compiler.NamespaceType.CelType()),
		cel.Variable(compiler.ObjectKey, cel.DynType),
		cel.Variable(compiler.OldObjectKey, cel.DynType),
		cel.Variable(compiler.RequestKey, compiler.RequestType.CelType()),
		cel.Types(compiler.NamespaceType.CelType()),
		cel.Types(compiler.RequestType.CelType()),
		cel.Variable(compiler.VariablesKey, compiler.VariablesType),
	)

	envSetVersion := compiler.KyvernoVersion
	base := environment.MustBaseEnvSet(envSetVersion)
	baseEnv, err := base.Env(environment.StoredExpressions)
	if err != nil {
		return nil, err
	}

	variablesProvider := compiler.NewVariablesProvider(baseEnv.CELTypeProvider())
	declProvider := apiservercel.NewDeclTypeProvider(compiler.NamespaceType, compiler.RequestType)
	declOptions, err := declProvider.EnvOptions(variablesProvider)
	if err != nil {
		return nil, err
	}
	baseOpts = append(baseOpts, declOptions...)

	libEnvOpts := []cel.EnvOption{
		ext.NativeTypes(reflect.TypeFor[libs.Exception](), ext.ParseStructTags(true)),
		cel.Variable(compiler.ExceptionsKey, types.NewObjectType("libs.Exception")),
		globalcontext.Lib(
			globalcontext.Context{ContextInterface: libsctx},
			globalcontext.Latest(),
		),
		resource.Lib(
			resource.Context{ContextInterface: libsctx},
			namespace,
			resource.Latest(),
		),
		image.Lib(
			image.Latest(),
		),
		imagedata.Lib(
			imagedata.Context{ContextInterface: libsctx},
			imagedata.Latest(),
			nil,
		),
		user.Lib(
			user.Latest(),
		),
		hash.Lib(
			hash.Latest(),
		),
		math.Lib(
			math.Latest(),
		),
		json.Lib(
			&json.JsonImpl{},
			json.Latest(),
		),
		yaml.Lib(
			&yaml.YamlImpl{},
			yaml.Latest(),
		),
		random.Lib(
			random.Latest(),
		),
		x509.Lib(
			x509.Latest(),
		),
		time.Lib(
			time.Latest(),
		),
		transform.Lib(
			transform.Latest(),
		),
		gzip.Lib(
			gzip.Latest(),
		),
		http.Lib(
			http.Context{ContextInterface: libs.NewMockAwareHTTPContext(compiler.NewLazyCELHTTPContext(namespace), libsctx.GetHTTPMocks())},
			http.Latest(),
		),
	}

	extendedBase, err := base.Extend(
		environment.VersionedOptions{
			IntroducedVersion: envSetVersion,
			EnvOptions:        baseOpts,
		},
		environment.VersionedOptions{
			IntroducedVersion: envSetVersion,
			EnvOptions:        libEnvOpts,
		},
	)
	if err != nil {
		return nil, err
	}

	return extendedBase.Env(environment.StoredExpressions)
}
