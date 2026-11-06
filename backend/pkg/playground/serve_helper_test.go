package playground_test

import (
	"os"
	"testing"

	"github.com/kyverno/playground/backend/pkg/playground"
)

// requireServe skips tests that need a real evaluation backend unless KYVERNO_SERVE_URL points at
// a `KYVERNO_EXPERIMENTAL=true kyverno serve --allow-request-policies` instance.
func requireServe(t *testing.T) {
	t.Helper()
	url := os.Getenv("KYVERNO_SERVE_URL")
	if url == "" {
		t.Skip("set KYVERNO_SERVE_URL to a `kyverno serve --allow-request-policies` instance")
	}
	playground.SetServeURL(url)
}
