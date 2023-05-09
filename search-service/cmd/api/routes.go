package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/log-entry", LogSearchEntry)
	mux.Post("/search-entry", SearchOneEntry)
	mux.Post("/get-pdf", SearchPDF)

	return mux
}
