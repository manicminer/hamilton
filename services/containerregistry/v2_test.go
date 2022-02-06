package containerregistry

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestV2(t *testing.T) {
	fa := testNewFakeAuthorizer(t)
	h := testNewACRHandler(t)
	httpServer := httptest.NewTLSServer(h.handler(t))
	h.serverURL = httpServer.URL
	cr := NewContainerRegistryClient(fa, httpServer.URL, "")
	cr.WithHttpClient(httpServer.Client())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testV2Check(t, ctx, cr)
}

func testV2Check(t *testing.T, ctx context.Context, cr *ContainerRegistryClient) {
	t.Helper()

	v2Supported, err := cr.V2Check(ctx)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if !v2Supported {
		t.Fatalf("expected v2 to be supported")
	}
}

func (h *testACRHandler) v2Handler(t *testing.T, w http.ResponseWriter, r *http.Request) {
	t.Helper()

	err := h.validateV2CheckRequest(t, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) //nolint
		return
	}
}

func (h *testACRHandler) validateV2CheckRequest(t *testing.T, r *http.Request) error {
	t.Helper()

	if r.Method != http.MethodGet {
		return fmt.Errorf("expected method to be GET, received: %s", r.Method)
	}

	path := r.URL.Path
	expectedPath := "/v2/"
	if path != expectedPath {
		return fmt.Errorf("expected path %q, received path: %s", expectedPath, path)
	}

	if h.fakeError != nil {
		return h.fakeError
	}

	return nil
}
