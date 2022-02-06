package containerregistry

import (
	"net/http/httptest"
	"testing"
)

func TestNewAuthorizer(t *testing.T) {
	fa := testNewFakeAuthorizer(t)
	h := testNewACRHandler(t)
	httpServer := httptest.NewTLSServer(h.handler(t))
	h.serverURL = httpServer.URL
	cr := NewContainerRegistryClient(fa, httpServer.URL, "")
	cr.WithHttpClient(httpServer.Client())

	testNewAuthorizer(t, cr)
}

func testNewAuthorizer(t *testing.T, cr *ContainerRegistryClient) {
	t.Helper()

	authorizer, err := cr.NewAuthorizer()
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	token, err := authorizer.Token()
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if !token.Valid() {
		t.Fatalf("expected token to be valid")
	}
}
