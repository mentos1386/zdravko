package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"code.tjo.space/mentos1386/zdravko/internal"
	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/handlers"
	"code.tjo.space/mentos1386/zdravko/web/static"
)

func main() {
	cfg := config.NewConfig()

	r := mux.NewRouter()

	db, query, err := internal.ConnectToDatabase(cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")

	h := handlers.NewBaseHandler(db, query, cfg)

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
	r.HandleFunc("/settings/healthchecks/create", h.Authenticated(h.SettingsHealthchecksCreateGET)).Methods("GET")
	r.HandleFunc("/settings/healthchecks/create", h.Authenticated(h.SettingsHealthchecksCreatePOST)).Methods("POST")
	r.HandleFunc("/settings/healthchecks/{id}", h.Authenticated(h.SettingsHealthchecksDescribeGET)).Methods("GET")

	// OAuth2
	r.HandleFunc("/oauth2/login", h.OAuth2LoginGET).Methods("GET")
	r.HandleFunc("/oauth2/callback", h.OAuth2CallbackGET).Methods("GET")
	r.HandleFunc("/oauth2/logout", h.Authenticated(h.OAuth2LogoutGET)).Methods("GET")

	// Temporal UI
	r.PathPrefix("/temporal").HandlerFunc(h.Authenticated(h.Temporal))

	// 404
	r.PathPrefix("/").HandlerFunc(h.Error404).Methods("GET")

	log.Println("Server started on", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
