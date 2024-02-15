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

	db, query, err := internal.ConnectToDatabase(config.DatabasePath)
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
	r.HandleFunc("/settings", h.Authenticated(h.SettingsOverviewGET)).Methods("GET")
	r.HandleFunc("/settings/healthchecks", h.Authenticated(h.SettingsHealthchecksGET)).Methods("GET")

	// OAuth2
	r.HandleFunc("/oauth2/login", h.OAuth2LoginGET).Methods("GET")
	r.HandleFunc("/oauth2/callback", h.OAuth2CallbackGET).Methods("GET")
	r.HandleFunc("/oauth2/logout", h.Authenticated(h.OAuth2LogoutGET)).Methods("GET")

	// Temporal UI
	r.PathPrefix("/temporal").HandlerFunc(h.Authenticated(h.Temporal))

	// 404
	r.PathPrefix("/").HandlerFunc(h.Error404).Methods("GET")

	log.Println("Server started on", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, r))
}
