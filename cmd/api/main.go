package main

import (
	"net/http"

	"github.com/ffernan01/animal-api/cmd/api/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	handler.Configure(r)

	http.ListenAndServe(":8080", r)
}
