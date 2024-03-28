package msgraph

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestClient_GetWithError(t *testing.T) {
	// This creates a listener on a random available port.
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()

	hc := NewClient(VersionBeta)
	hc.Endpoint = fmt.Sprintf("https://localhost:%d/", port)
	hc.RetryableClient.RetryMax = 2

	_, _, _, err = hc.Get(context.Background(), GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity: "/users/test",
		},
	})
	if err == nil {
		t.Error("expected to get an error, got nil")
	}
	if msg := err.Error(); !strings.Contains(msg, "connect: connection refused") {
		log.Fatalf("got %s, want message with 'connection refused'", msg)
	}
}

func TestClient_GetWithResponseAndError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/beta/users/test" {
			n, _ := strconv.Atoi(r.FormValue("n"))
			if n < 15 {
				http.Redirect(w, r, fmt.Sprintf("%s?n=%d", r.URL.Path, 1), http.StatusTemporaryRedirect)
				return
			}
		}
	}))
	defer ts.Close()

	hc := NewClient(VersionBeta)
	hc.Endpoint = ts.URL
	hc.RetryableClient.RetryMax = 2

	_, _, _, err := hc.Get(context.Background(), GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity: "/users/test",
		},
	})
	if err == nil {
		t.Error("expected to get an error, got nil")
	}
	if msg := err.Error(); !strings.Contains(msg, "stopped after 10 redirects") {
		log.Fatalf("got %s, want message with 'stopped after 10 redirects'", msg)
	}
}
