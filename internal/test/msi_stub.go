package test

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func MsiStubServer(ctx context.Context, port int, token string) chan bool {
	handler := http.NewServeMux()

	handler.HandleFunc("/metadata/identity/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(w, `{"access_token":"%s","client_id":"00000000-0000-0000-0000-000000000000","expires_in":"86391","expires_on":"1611701390","ext_expires_in":"86399","not_before":"1611614690","resource":"https://graph.microsoft.com/","token_type":"Bearer"}`, token)
	})

	handler.HandleFunc("/metadata", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, `
`)
	})

	done := make(chan bool, 1)
	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: handler,
	}

	go func() {
		<-done
		err := server.Shutdown(ctx)
		if err != nil {
			log.Fatal("failed to gracefully shut down MSI stub server")
		}
	}()

	go func() {
		log.Println("MSI Stub Service listening on 127.0.0.1:8080")
		log.Fatal(server.ListenAndServe())
	}()

	return done
}
