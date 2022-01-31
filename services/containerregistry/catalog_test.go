package containerregistry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCatalogClient(t *testing.T) {
	fa := testNewFakeAuthorizer(t)
	h := testNewACRHandler(t)
	httpServer := httptest.NewTLSServer(h.handler(t))
	h.serverURL = httpServer.URL
	cr := NewContainerRegistryClient(fa, httpServer.URL, "")
	cr.WithHttpClient(httpServer.Client())
	catalogClient := NewCatalogClient(cr)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testCatalogClientList(t, ctx, catalogClient)
	testCatalogClientUpdateAttributes(t, ctx, catalogClient)
	testCatalogClientGetAttributes(t, ctx, catalogClient)
	testCatalogClientDelete(t, ctx, catalogClient)
}

func testCatalogClientList(t *testing.T, ctx context.Context, catalogClient *CatalogClient) {
	t.Helper()

	repositories, err := catalogClient.List(ctx, nil, nil)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if len(repositories) != 3 {
		t.Fatalf("expected three repositories")
	}
}

func testCatalogClientGetAttributes(t *testing.T, ctx context.Context, catalogClient *CatalogClient) {
	t.Helper()

	imageName := "foobar"
	res, err := catalogClient.GetAttributes(ctx, imageName)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if res.ImageName != imageName {
		t.Fatalf("expected to receive image name %q, but got: %s", imageName, res.ImageName)
	}
}

func testCatalogClientUpdateAttributes(t *testing.T, ctx context.Context, catalogClient *CatalogClient) {
	t.Helper()

	toBoolPtr := func(v bool) *bool { return &v }
	err := catalogClient.UpdateAttributes(ctx, "foobar", RepositoryChangeableAttributes{ListEnabled: toBoolPtr(true)})
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}
}

func testCatalogClientDelete(t *testing.T, ctx context.Context, catalogClient *CatalogClient) {
	t.Helper()

	res, err := catalogClient.Delete(ctx, "foobar")
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if res.ManifestsDeleted[0] != "sha256:e31831d63f77a0a6d74ef5b16df619a50808dac842190d07ae24e8b520d159fa" {
		t.Fatal("expected to receive deleted manifest 'sha256:e31831d63f77a0a6d74ef5b16df619a50808dac842190d07ae24e8b520d159fa'")
	}

	if res.TagsDeleted[0] != "latest" {
		t.Fatal("expected tor eceive deleted tag 'latest'")
	}
}

func (h *testACRHandler) catalogHandler(t *testing.T, w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(r.URL.Path, "/", 4)
	imageName := parts[3]
	if imageName == "_catalog" {
		h.catalogListHandler(t, w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.catalogGetAttributesHandler(t, w, r, imageName)
	case http.MethodPatch:
		h.catalogUpdateAttributesHandler(t, w, r, imageName)
	case http.MethodDelete:
		h.catalogDeleteHandler(t, w, r, imageName)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unknown method: %s", r.Method)))
		return
	}
}

func (h *testACRHandler) catalogListHandler(t *testing.T, w http.ResponseWriter, r *http.Request) {
	err := h.validateCatalogListRequest(t, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response := struct {
		Repositories []string `json:"repositories"`
	}{
		Repositories: []string{"foo", "bar", "baz"},
	}

	json.NewEncoder(w).Encode(response)
}

func (h *testACRHandler) validateCatalogListRequest(t *testing.T, r *http.Request) error {
	t.Helper()

	if r.Method != http.MethodGet {
		return fmt.Errorf("expected method to be GET, received: %s", r.Method)
	}

	path := r.URL.Path
	if path != "/acr/v1/_catalog" {
		return fmt.Errorf("expected path '/acr/v1/_catalog', received path: %s", path)
	}

	query := r.URL.Query()
	if len(query) > 2 {
		return fmt.Errorf("expected query to contain max of 2 parameters, received: %s", query)
	}

	if h.fakeError != nil {
		return h.fakeError
	}

	return nil
}

func (h *testACRHandler) catalogGetAttributesHandler(t *testing.T, w http.ResponseWriter, r *http.Request, imageName string) {
	err := h.validateCatalogGetAttributesRequest(t, r, imageName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response := RepositoryAttributesResponse{
		Registry:       "foobar.azurecr.io",
		ImageName:      imageName,
		CreatedTime:    time.Now().Add(-60 * time.Minute),
		LastUpdateTime: time.Now().Add(-30 * time.Minute),
		ManifestCount:  1,
		TagCount:       1,
		ChangeableAttributes: RepositoryChangeableAttributesResponse{
			DeleteEnabled:   true,
			WriteEnabled:    true,
			ReadEnabled:     true,
			ListEnabled:     true,
			TeleportEnabled: false,
		},
	}

	json.NewEncoder(w).Encode(response)
}

func (h *testACRHandler) validateCatalogGetAttributesRequest(t *testing.T, r *http.Request, imageName string) error {
	t.Helper()

	if r.Method != http.MethodGet {
		return fmt.Errorf("expected method to be GET, received: %s", r.Method)
	}

	path := r.URL.Path
	expectedPath := fmt.Sprintf("/acr/v1/%s", imageName)
	if path != expectedPath {
		return fmt.Errorf("expected path %q, received path: %s", expectedPath, path)
	}

	query := r.URL.Query()
	if len(query) > 2 {
		return fmt.Errorf("expected query to contain max of 2 parameters, received: %s", query)
	}

	if h.fakeError != nil {
		return h.fakeError
	}

	return nil
}

func (h *testACRHandler) catalogUpdateAttributesHandler(t *testing.T, w http.ResponseWriter, r *http.Request, imageName string) {
	err := h.validateCatalogUpdateAttributesRequest(t, r, imageName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (h *testACRHandler) validateCatalogUpdateAttributesRequest(t *testing.T, r *http.Request, imageName string) error {
	t.Helper()

	if r.Method != http.MethodPatch {
		return fmt.Errorf("expected method to be PATCH, received: %s", r.Method)
	}

	path := r.URL.Path
	expectedPath := fmt.Sprintf("/acr/v1/%s", imageName)
	if path != expectedPath {
		return fmt.Errorf("expected path %q, received path: %s", expectedPath, path)
	}

	if h.fakeError != nil {
		return h.fakeError
	}

	return nil
}

func (h *testACRHandler) catalogDeleteHandler(t *testing.T, w http.ResponseWriter, r *http.Request, imageName string) {
	err := h.validateCatalogDeleteRequest(t, r, imageName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	response := RepositoryDeleteResponse{
		ManifestsDeleted: []string{"sha256:e31831d63f77a0a6d74ef5b16df619a50808dac842190d07ae24e8b520d159fa"},
		TagsDeleted:      []string{"latest"},
	}

	json.NewEncoder(w).Encode(response)
}

func (h *testACRHandler) validateCatalogDeleteRequest(t *testing.T, r *http.Request, imageName string) error {
	t.Helper()

	if r.Method != http.MethodDelete {
		return fmt.Errorf("expected method to be DELETE, received: %s", r.Method)
	}

	path := r.URL.Path
	expectedPath := fmt.Sprintf("/acr/v1/%s", imageName)
	if path != expectedPath {
		return fmt.Errorf("expected path %q, received path: %s", expectedPath, path)
	}

	if h.fakeError != nil {
		return h.fakeError
	}

	return nil
}
