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

func TestTag(t *testing.T) {
	fa := testNewFakeAuthorizer(t)
	h := testNewACRHandler(t)
	httpServer := httptest.NewTLSServer(h.handler(t))
	h.serverURL = httpServer.URL
	cr := NewContainerRegistryClient(fa, httpServer.URL, "")
	cr.WithHttpClient(httpServer.Client())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testTagList(t, ctx, cr)
	testTagUpdateAttributes(t, ctx, cr)
	testTagGetAttributes(t, ctx, cr)
	testTagDelete(t, ctx, cr)
}

func testTagList(t *testing.T, ctx context.Context, cr *ContainerRegistryClient) {
	t.Helper()

	imageName := "foobar"
	tagList, err := cr.TagList(ctx, imageName, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if len(tagList.Tags) != 3 {
		t.Fatalf("expected three tags")
	}
}

func testTagGetAttributes(t *testing.T, ctx context.Context, cr *ContainerRegistryClient) {
	t.Helper()

	imageName := "foobar"
	reference := "latest"
	res, err := cr.TagGetAttributes(ctx, imageName, reference)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if res.ImageName != imageName {
		t.Fatalf("expected to receive image name %q, but got: %s", imageName, res.ImageName)
	}

	if res.Tag.Name != reference {
		t.Fatalf("expected to receive tag name %q, but got: %s", reference, res.Tag.Name)
	}
}

func testTagUpdateAttributes(t *testing.T, ctx context.Context, cr *ContainerRegistryClient) {
	t.Helper()

	toBoolPtr := func(v bool) *bool { return &v }
	err := cr.TagUpdateAttributes(ctx, "foobar", "latest", TagChangeableAttributes{ListEnabled: toBoolPtr(true)})
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}
}

func testTagDelete(t *testing.T, ctx context.Context, cr *ContainerRegistryClient) {
	t.Helper()

	err := cr.TagDelete(ctx, "foobar", "latest")
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}
}

func (h *testACRHandler) tagHandler(t *testing.T, w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(r.URL.Path, "/", 4)
	tagPath := parts[3]
	if strings.HasSuffix(tagPath, "/_tags") {
		imageName := strings.TrimSuffix(tagPath, "/_tags")
		h.tagListHandler(t, w, r, imageName)
		return
	}

	tagParts := strings.SplitN(r.URL.Path, "/_tags/", 2)
	if len(tagParts) != 2 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("expected tagParts to be of length 2: %s", tagParts))) //nolint
		return
	}

	imageName := tagParts[0]
	reference := tagParts[1]

	switch r.Method {
	case http.MethodGet:
		h.tagGetAttributesHandler(t, w, r, imageName, reference)
	case http.MethodPatch:
		h.tagUpdateAttributesHandler(t, w, r, imageName, reference)
	case http.MethodDelete:
		h.tagDeleteHandler(t, w, r, imageName, reference)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unknown method: %s", r.Method))) //nolint
		return
	}
}

func (h *testACRHandler) tagListHandler(t *testing.T, w http.ResponseWriter, r *http.Request, imageName string) {
	err := h.validateTagListRequest(t, r, imageName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) //nolint
		return
	}

	response := struct {
		Repositories []string `json:"repositories"`
	}{
		Repositories: []string{"foo", "bar", "baz"},
	}

	json.NewEncoder(w).Encode(response) //nolint
}

func (h *testACRHandler) validateTagListRequest(t *testing.T, r *http.Request, imageName string) error {
	t.Helper()

	if r.Method != http.MethodGet {
		return fmt.Errorf("expected method to be GET, received: %s", r.Method)
	}

	path := r.URL.Path
	expectedPath := fmt.Sprintf("/acr/v1/%s/_tags", imageName)
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

func (h *testACRHandler) tagGetAttributesHandler(t *testing.T, w http.ResponseWriter, r *http.Request, imageName string, reference string) {
	err := h.validateTagGetAttributesRequest(t, r, imageName, reference)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) //nolint
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

	json.NewEncoder(w).Encode(response) //nolint
}

func (h *testACRHandler) validateTagGetAttributesRequest(t *testing.T, r *http.Request, imageName string, reference string) error {
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

func (h *testACRHandler) tagUpdateAttributesHandler(t *testing.T, w http.ResponseWriter, r *http.Request, imageName string, reference string) {
	err := h.validateTagUpdateAttributesRequest(t, r, imageName, reference)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) //nolint
		return
	}
}

func (h *testACRHandler) validateTagUpdateAttributesRequest(t *testing.T, r *http.Request, imageName string, reference string) error {
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

func (h *testACRHandler) tagDeleteHandler(t *testing.T, w http.ResponseWriter, r *http.Request, imageName string, reference string) {
	err := h.validateTagDeleteRequest(t, r, imageName, reference)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) //nolint
		return
	}

	w.WriteHeader(http.StatusAccepted)
	response := RepositoryDeleteResponse{
		ManifestsDeleted: []string{"sha256:e31831d63f77a0a6d74ef5b16df619a50808dac842190d07ae24e8b520d159fa"},
		TagsDeleted:      []string{"latest"},
	}

	json.NewEncoder(w).Encode(response) //nolint
}

func (h *testACRHandler) validateTagDeleteRequest(t *testing.T, r *http.Request, imageName string, reference string) error {
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
