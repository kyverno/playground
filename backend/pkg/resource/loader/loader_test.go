package loader

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/client-go/openapi"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
	"sigs.k8s.io/kubectl-validate/pkg/validatorfactory"

	"github.com/kyverno/playground/backend/data"
)

type errClient struct{}

func (_ errClient) Paths() (map[string]openapi.GroupVersion, error) {
	return nil, errors.New("error")
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		client  openapi.Client
		want    Loader
		wantErr bool
	}{{
		name:    "err client",
		client:  errClient{},
		wantErr: true,
	}, {
		name:   "builtin",
		client: openapiclient.NewHardcodedBuiltins("1.27"),
		want: func() Loader {
			factory, err := validatorfactory.New(openapiclient.NewHardcodedBuiltins("1.27"))
			require.NoError(t, err)
			return &loader{
				factory: factory,
			}
		}(),
	}, {
		name:    "invalid local",
		client:  openapiclient.NewLocalFiles(data.Schemas(), "blam"),
		wantErr: true,
	}, {
		name:   "composite - no clients",
		client: openapiclient.NewComposite(),
		want: func() Loader {
			factory, err := validatorfactory.New(openapiclient.NewComposite())
			require.NoError(t, err)
			return &loader{
				factory: factory,
			}
		}(),
	}, {
		name:   "composite - err client",
		client: openapiclient.NewComposite(errClient{}),
		want: func() Loader {
			factory, err := validatorfactory.New(openapiclient.NewComposite(errClient{}))
			require.NoError(t, err)
			return &loader{
				factory: factory,
			}
		}(),
	}, {
		name:   "composite - with err client",
		client: openapiclient.NewComposite(openapiclient.NewHardcodedBuiltins("1.27"), errClient{}),
		want: func() Loader {
			factory, err := validatorfactory.New(openapiclient.NewComposite(openapiclient.NewHardcodedBuiltins("1.27"), errClient{}))
			require.NoError(t, err)
			return &loader{
				factory: factory,
			}
		}(),
	}, {
		name:   "composite - invalid local",
		client: openapiclient.NewComposite(openapiclient.NewLocalFiles(data.Schemas(), "blam")),
		want: func() Loader {
			factory, err := validatorfactory.New(openapiclient.NewComposite(openapiclient.NewLocalFiles(data.Schemas(), "blam")))
			require.NoError(t, err)
			return &loader{
				factory: factory,
			}
		}(),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
