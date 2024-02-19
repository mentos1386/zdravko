package handlers

import (
	"encoding/json"
	"net/http"
)

type ApiV1WorkersConnectGETResponse struct {
	Endpoint string `json:"endpoint"`
	Group    string `json:"group"`
	Slug     string `json:"slug"`
}

func (h *BaseHandler) ApiV1WorkersConnectGET(w http.ResponseWriter, r *http.Request, principal *AuthenticatedPrincipal) {
	// Json response containing temporal endpoint
	w.Header().Set("Content-Type", "application/json")

	response := ApiV1WorkersConnectGETResponse{
		Endpoint: h.config.Temporal.ServerHost,
		Group:    principal.Worker.Group,
		Slug:     principal.Worker.Slug,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(responseJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
