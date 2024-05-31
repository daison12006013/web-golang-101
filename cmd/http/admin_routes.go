package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func AdminRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Admin route"))
	})
	return r
}
