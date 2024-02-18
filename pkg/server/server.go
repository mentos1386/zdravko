package server

import (
	"context"
	"log"
	"net/http"

	"code.tjo.space/mentos1386/zdravko/internal"
	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/handlers"
	"code.tjo.space/mentos1386/zdravko/internal/temporal"
	"code.tjo.space/mentos1386/zdravko/web/static"
	"github.com/gorilla/mux"
)

type Server struct {
	server *http.Server
	cfg    *config.Config
}

func NewServer(cfg *config.Config) (*Server, error) {
	return &Server{
		cfg: cfg,
	}, nil
}

func (s *Server) Name() string {
	return "HTTP WEB and API Server"
}

func (s *Server) Start() error {
	r := mux.NewRouter()

	db, query, err := internal.ConnectToDatabase(s.cfg.DatabasePath)
	if err != nil {
		return err
	}
	log.Println("Connected to database")

	temporalClient, err := temporal.ConnectServerToTemporal(s.cfg)
	if err != nil {
		return err
	}
	log.Println("Connected to Temporal")

	h := handlers.NewBaseHandler(db, query, temporalClient, s.cfg)

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
	r.HandleFunc("/settings/healthchecks/{slug}", h.Authenticated(h.SettingsHealthchecksDescribeGET)).Methods("GET")
	r.HandleFunc("/settings/workers", h.Authenticated(h.SettingsWorkersGET)).Methods("GET")
	r.HandleFunc("/settings/workers/create", h.Authenticated(h.SettingsWorkersCreateGET)).Methods("GET")
	r.HandleFunc("/settings/workers/create", h.Authenticated(h.SettingsWorkersCreatePOST)).Methods("POST")
	r.HandleFunc("/settings/workers/{slug}", h.Authenticated(h.SettingsWorkersDescribeGET)).Methods("GET")
	r.HandleFunc("/settings/workers/{slug}/token", h.Authenticated(h.SettingsWorkersTokenGET)).Methods("GET")

	// OAuth2
	r.HandleFunc("/oauth2/login", h.OAuth2LoginGET).Methods("GET")
	r.HandleFunc("/oauth2/callback", h.OAuth2CallbackGET).Methods("GET")
	r.HandleFunc("/oauth2/logout", h.Authenticated(h.OAuth2LogoutGET)).Methods("GET")

	// Temporal UI
	r.PathPrefix("/temporal").HandlerFunc(h.Authenticated(h.Temporal))

	// 404
	r.PathPrefix("/").HandlerFunc(h.Error404).Methods("GET")

	s.server = &http.Server{
		Addr:    ":" + s.cfg.Port,
		Handler: r,
	}

	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx := context.Background()
	return s.server.Shutdown(ctx)
}
