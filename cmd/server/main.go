package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"code.tjo.space/mentos1386/zdravko/internal"
	"code.tjo.space/mentos1386/zdravko/internal/pages"
)

func main() {
	r := mux.NewRouter()

	// Server static files
	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(internal.Static)))

	r.HandleFunc("/", pages.Index).Methods("GET")

	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
