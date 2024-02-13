package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func (h *BaseHandler) Temporal(w http.ResponseWriter, r *http.Request, user *AuthenticatedUser) {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Host:   h.config.TEMPORAL_UI_HOST,
		Scheme: "http",
	})

	// TODO: Maybe add a "navbar" in html to go back to Zdravko?
	proxy.ModifyResponse = func(response *http.Response) error {
		// Read and update the response here

		// The response here is response from server (proxy B if this is at proxy A)
		// It is a pointer, so can be modified to update in place
		// It will not be called if Proxy B is unreachable
		return nil
	}

	proxy.ServeHTTP(w, r)
}
