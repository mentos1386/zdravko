package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"code.tjo.space/mentos1386/zdravko/internal"
	"code.tjo.space/mentos1386/zdravko/internal/pages"
	"code.tjo.space/mentos1386/zdravko/internal/static"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	r := mux.NewRouter()

	db, err := internal.ConnectToDatabase()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")

	page := pages.NewPageHandler(db)

	// Server static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(static.Static))))

	r.HandleFunc("/", page.Index).Methods("GET")
	r.HandleFunc("/settings", page.Settings).Methods("GET")

	log.Println("Server started on", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
