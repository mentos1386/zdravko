package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"code.tjo.space/mentos1386/zdravko/internal"
	"code.tjo.space/mentos1386/zdravko/internal/handlers"
	"code.tjo.space/mentos1386/zdravko/web/static"
)

func main() {
	config := internal.NewConfig()

	r := mux.NewRouter()

	db, query, err := internal.ConnectToDatabase(config.SQLITE_DB_PATH)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")

	h := handlers.NewBaseHandler(db, query, config)

	// Health
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		d, err := db.DB()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = d.Ping()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		_, err = w.Write([]byte("OK"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Server static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(static.Static))))

	r.HandleFunc("/", h.Index).Methods("GET")

	// Authenticated routes
	r.HandleFunc("/settings", h.Authenticated(h.Settings)).Methods("GET")

	// OAuth2
	r.HandleFunc("/oauth2/login", h.OAuth2LoginGET).Methods("GET")
	r.HandleFunc("/oauth2/callback", h.OAuth2CallbackGET).Methods("GET")
	r.HandleFunc("/oauth2/logout", h.Authenticated(h.OAuth2LogoutGET)).Methods("GET")

	log.Println("Server started on", config.PORT)
	log.Fatal(http.ListenAndServe(":"+config.PORT, r))
}
