package containerregistry

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2"
)

func TestContainerRegistryClient(t *testing.T) {
	fa := &testFakeAuthorizer{}
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}
	httpServer := httptest.NewServer(http.HandlerFunc(handler))
	cr := NewContainerRegistryClient(fa, httpServer.URL, "")
	ctx := context.Background()
	token, err := cr.ExchangeToken(ctx)
	if err != nil {
		t.FailNow()
	}

	t.Fatalf("containue here: %s", token)
}

type testFakeAuthorizer struct{}

func (fa *testFakeAuthorizer) Token() (*oauth2.Token, error) {
	return &oauth2.Token{}, nil
}

func (fa *testFakeAuthorizer) AuxiliaryTokens() ([]*oauth2.Token, error) {
	return nil, nil
}
