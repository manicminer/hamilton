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

	if tagList.ImageName != imageName {
		t.Fatalf("expected to receive image name %q, but got: %s", imageName, tagList.ImageName)
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

	tagParts := strings.SplitN(tagPath, "/_tags/", 2)
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

	response := TagList{
		Registry:  "foo.azureacr.io",
		ImageName: imageName,
		Tags: []Tag{
			{
				Name: "tag1",
			},
			{
				Name: "tag2",
			},
			{
				Name: "tag3",
			},
		},
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

	response := TagAttributesResponse{
		Registry:  "foo.azurecr.io",
		ImageName: imageName,
		Tag: Tag{
			Name: reference,
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
	expectedPath := fmt.Sprintf("/acr/v1/%s/_tags/%s", imageName, reference)
	if path != expectedPath {
		return fmt.Errorf("expected path %q, received path: %s", expectedPath, path)
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
	expectedPath := fmt.Sprintf("/acr/v1/%s/_tags/%s", imageName, reference)
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
}

func (h *testACRHandler) validateTagDeleteRequest(t *testing.T, r *http.Request, imageName string, reference string) error {
	t.Helper()

	if r.Method != http.MethodDelete {
		return fmt.Errorf("expected method to be DELETE, received: %s", r.Method)
	}

	path := r.URL.Path
	expectedPath := fmt.Sprintf("/acr/v1/%s/_tags/%s", imageName, reference)
	if path != expectedPath {
		return fmt.Errorf("expected path %q, received path: %s", expectedPath, path)
	}

	if h.fakeError != nil {
		return h.fakeError
	}

	return nil
}
