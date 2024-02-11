package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"code.tjo.space/mentos1386/zdravko/internal"
	"code.tjo.space/mentos1386/zdravko/internal/pages"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	r := mux.NewRouter()

	// Server static files
	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(internal.Static)))

	r.HandleFunc("/", pages.Index).Methods("GET")
	r.HandleFunc("/settings", pages.Settings).Methods("GET")

	log.Println("Server started on", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
