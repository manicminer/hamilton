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

func TestManifest(t *testing.T) {
	fa := testNewFakeAuthorizer(t)
	h := testNewACRHandler(t)
	httpServer := httptest.NewTLSServer(h.handler(t))
	h.serverURL = httpServer.URL
	cr := NewContainerRegistryClient(fa, httpServer.URL, "")
	cr.WithHttpClient(httpServer.Client())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testManifestList(t, ctx, cr)
	testManifestGet(t, ctx, cr)
	testManifestUpdateAttributes(t, ctx, cr)
	testManifestGetAttributes(t, ctx, cr)
	testManifestDelete(t, ctx, cr)
}

func testManifestList(t *testing.T, ctx context.Context, cr *ContainerRegistryClient) {
	t.Helper()

	imageName := "foobar"
	manifestList, err := cr.ManifestList(ctx, imageName, nil, nil, nil)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if len(manifestList.Manifests) != 3 {
		t.Fatalf("expected three manifests")
	}

	if manifestList.ImageName != imageName {
		t.Fatalf("expected to receive image name %q, but got: %s", imageName, manifestList.ImageName)
	}
}

func testManifestGet(t *testing.T, ctx context.Context, cr *ContainerRegistryClient) {
	t.Helper()

	imageName := "foobar"
	reference := "latest"
	manifest, err := cr.ManifestGet(ctx, imageName, reference)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if len(manifest.Layers) != 5 {
		t.Fatalf("expected five layers")
	}
}

func testManifestGetAttributes(t *testing.T, ctx context.Context, cr *ContainerRegistryClient) {
	t.Helper()

	imageName := "foobar"
	reference := "latest"
	res, err := cr.ManifestGetAttributes(ctx, imageName, reference)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	if res.ImageName != imageName {
		t.Fatalf("expected to receive image name %q, but got: %s", imageName, res.ImageName)
	}
}

func testManifestUpdateAttributes(t *testing.T, ctx context.Context, cr *ContainerRegistryClient) {
	t.Helper()

	toBoolPtr := func(v bool) *bool { return &v }
	err := cr.ManifestUpdateAttributes(ctx, "foobar", "latest", ManifestChangeableAttributes{ListEnabled: toBoolPtr(true)})
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}
}

func testManifestDelete(t *testing.T, ctx context.Context, cr *ContainerRegistryClient) {
	t.Helper()

	err := cr.ManifestDelete(ctx, "foobar", "sha256:0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}
}

func (h *testACRHandler) manifestHandler(t *testing.T, w http.ResponseWriter, r *http.Request) {
	t.Helper()

	apiVersion := 1
	if strings.HasPrefix(r.URL.Path, "/v2/") {
		apiVersion = 2
	}

	switch apiVersion {
	case 1:
		h.manifestV1Handler(t, w, r)
	case 2:
		h.manifestV2Handler(t, w, r)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("you shouldn't be able to reach this error..."))) //nolint
		return
	}
}

func (h *testACRHandler) manifestV1Handler(t *testing.T, w http.ResponseWriter, r *http.Request) {
	t.Helper()

	parts := strings.SplitN(r.URL.Path, "/", 4)
	manifestPath := parts[3]
	if strings.HasSuffix(manifestPath, "/_manifests") {
		imageName := strings.TrimSuffix(manifestPath, "/_manifests")
		h.manifestListHandler(t, w, r, imageName)
		return
	}

	manifestParts := strings.SplitN(manifestPath, "/_manifests/", 2)
	if len(manifestParts) != 2 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("expected manifestParts to be of length 2: %s", manifestParts))) //nolint
		return
	}

	imageName := manifestParts[0]
	reference := manifestParts[1]

	switch r.Method {
	case http.MethodGet:
		h.manifestGetAttributesHandler(t, w, r, imageName, reference)
	case http.MethodPatch:
		h.manifestUpdateAttributesHandler(t, w, r, imageName, reference)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unknown method: %s", r.Method))) //nolint
		return
	}
}

func (h *testACRHandler) manifestV2Handler(t *testing.T, w http.ResponseWriter, r *http.Request) {
	t.Helper()

	// /v2/%s/manifests/%s
	manifestPath := strings.TrimPrefix(r.URL.Path, "/v2/")
	manifestParts := strings.SplitN(manifestPath, "/manifests/", 2)
	if len(manifestParts) != 2 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("expected manifestParts to be of length 2: %s", manifestParts))) //nolint
		return
	}

	imageName := manifestParts[0]
	reference := manifestParts[1]

	switch r.Method {
	case http.MethodGet:
		h.manifestGetHandler(t, w, r, imageName, reference)
	case http.MethodDelete:
		h.manifestDeleteHandler(t, w, r, imageName, reference)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unknown method: %s", r.Method))) //nolint
		return
	}
}

func (h *testACRHandler) manifestListHandler(t *testing.T, w http.ResponseWriter, r *http.Request, imageName string) {
	t.Helper()

	err := h.validateManifestListRequest(t, r, imageName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) //nolint
		return
	}

	response := ManifestList{
		Registry:  "foo.azurecr.io",
		ImageName: imageName,
		Manifests: []Manifest{
			{
				Digest: "sha256:0000000000000000000000000000000000000000000000000000000000000000",
			},
			{
				Digest: "sha256:0000000000000000000000000000000000000000000000000000000000000001",
			},
			{
				Digest: "sha256:0000000000000000000000000000000000000000000000000000000000000002",
			},
		},
	}

	json.NewEncoder(w).Encode(response) //nolint
}

func (h *testACRHandler) validateManifestListRequest(t *testing.T, r *http.Request, imageName string) error {
	t.Helper()

	if r.Method != http.MethodGet {
		return fmt.Errorf("expected method to be GET, received: %s", r.Method)
	}

	path := r.URL.Path
	expectedPath := fmt.Sprintf("/acr/v1/%s/_manifests", imageName)
	if path != expectedPath {
		return fmt.Errorf("expected path %q, received path: %s", expectedPath, path)
	}

	if h.fakeError != nil {
		return h.fakeError
	}

	return nil
}

func (h *testACRHandler) manifestGetHandler(t *testing.T, w http.ResponseWriter, r *http.Request, imageName string, reference string) {
	t.Helper()

	err := h.validateManifestGetRequest(t, r, imageName, reference)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) //nolint
		return
	}

	response := ManifestGetResponse{
		SchemaVersion: 2,
		MediaType:     "application/vnd.docker.distribution.manifest.v2+json",
		Config: ManifestConfigResponse{
			MediaType: "application/vnd.docker.container.image.v1+json",
			Size:      5824,
			Digest:    "sha256:691fbc2d44fff48357bba69ab0505b9bf12b2b250a925a84a0b8e8e7eed390b2",
		},
		Layers: []ManifestLayerResponse{
			{
				MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
				Size:      2014658,
				Digest:    "sha256:a073c86ecf9e0f29180e80e9638d4c741970695851ea48247276c32c57e40282",
			},
			{
				MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
				Size:      19778035,
				Digest:    "sha256:0e28711eb56d78f1e3dfde1807eba529d1346222bcd07d1cb1e436a18a0388bd",
			},
			{
				MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
				Size:      1074044,
				Digest:    "sha256:e460dd483fddb555911f7ed188c319fd97542c60e36843dcb1c5d753f733e1fa",
			},
			{
				MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
				Size:      5827,
				Digest:    "sha256:6aa301222093bfb8cf424ccb387f59e2c9510c3a30cca7fbcf8c954f88e6600c",
			},
			{
				MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
				Size:      568,
				Digest:    "sha256:9c5d80083a57d565f684e0155707204d497a5ad965279f92927452f15dae17e6",
			},
		},
	}

	json.NewEncoder(w).Encode(response) //nolint
}

func (h *testACRHandler) validateManifestGetRequest(t *testing.T, r *http.Request, imageName string, reference string) error {
	t.Helper()

	if r.Method != http.MethodGet {
		return fmt.Errorf("expected method to be GET, received: %s", r.Method)
	}

	path := r.URL.Path
	expectedPath := fmt.Sprintf("/v2/%s/manifests/%s", imageName, reference)
	if path != expectedPath {
		return fmt.Errorf("expected path %q, received path: %s", expectedPath, path)
	}

	if h.fakeError != nil {
		return h.fakeError
	}

	return nil
}

func (h *testACRHandler) manifestGetAttributesHandler(t *testing.T, w http.ResponseWriter, r *http.Request, imageName string, reference string) {
	t.Helper()

	err := h.validateManifestGetAttributesRequest(t, r, imageName, reference)
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

func (h *testACRHandler) validateManifestGetAttributesRequest(t *testing.T, r *http.Request, imageName string, reference string) error {
	t.Helper()

	if r.Method != http.MethodGet {
		return fmt.Errorf("expected method to be GET, received: %s", r.Method)
	}

	path := r.URL.Path
	expectedPath := fmt.Sprintf("/acr/v1/%s/_manifests/%s", imageName, reference)
	if path != expectedPath {
		return fmt.Errorf("expected path %q, received path: %s", expectedPath, path)
	}

	if h.fakeError != nil {
		return h.fakeError
	}

	return nil
}

func (h *testACRHandler) manifestUpdateAttributesHandler(t *testing.T, w http.ResponseWriter, r *http.Request, imageName string, reference string) {
	t.Helper()

	err := h.validateManifestUpdateAttributesRequest(t, r, imageName, reference)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) //nolint
		return
	}
}

func (h *testACRHandler) validateManifestUpdateAttributesRequest(t *testing.T, r *http.Request, imageName string, reference string) error {
	t.Helper()

	if r.Method != http.MethodPatch {
		return fmt.Errorf("expected method to be PATCH, received: %s", r.Method)
	}

	path := r.URL.Path
	expectedPath := fmt.Sprintf("/acr/v1/%s/_manifests/%s", imageName, reference)
	if path != expectedPath {
		return fmt.Errorf("expected path %q, received path: %s", expectedPath, path)
	}

	if h.fakeError != nil {
		return h.fakeError
	}

	return nil
}

func (h *testACRHandler) manifestDeleteHandler(t *testing.T, w http.ResponseWriter, r *http.Request, imageName string, reference string) {
	t.Helper()

	err := h.validateManifestDeleteRequest(t, r, imageName, reference)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) //nolint
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *testACRHandler) validateManifestDeleteRequest(t *testing.T, r *http.Request, imageName string, reference string) error {
	t.Helper()

	if r.Method != http.MethodDelete {
		return fmt.Errorf("expected method to be DELETE, received: %s", r.Method)
	}

	path := r.URL.Path
	expectedPath := fmt.Sprintf("/v2/%s/manifests/%s", imageName, reference)
	if path != expectedPath {
		return fmt.Errorf("expected path %q, received path: %s", expectedPath, path)
	}

	if h.fakeError != nil {
		return h.fakeError
	}

	return nil
}
