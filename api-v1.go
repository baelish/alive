package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/oapi-codegen/pkg/codegen"
)

// Your handler struct — implement all ServerInterface methods on this.
type MyAPI struct{}

func (a *MyAPI) GetBox(w http.ResponseWriter, r *http.Request, id string) {
	// your handler code here
}

// implement other methods...

func runAPI(ctx context.Context) {
	swagger, err := codegen.GetSwagger()
	if err != nil {
		log.Fatalf("failed to get swagger: %v", err)
	}

	router := chi.NewRouter()
	router.Use(codegen.OapiRequestValidator(swagger))

	// Register your implementation
	api := MyAPI{}
	codegen.HandlerFromMux(api, router)

	port := ":8080" // or from config
	log.Printf("Listening on %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
