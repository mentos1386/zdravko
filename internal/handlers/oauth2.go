package handlers

import (
	"context"
	"io"
	"net/http"

	"code.tjo.space/mentos1386/zdravko/internal"
	"golang.org/x/oauth2"
)

func newOAuth2(config *internal.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.OAUTH2_CLIENT_ID,
		ClientSecret: config.OAUTH2_CLIENT_SECRET,
		Scopes:       config.OAUTH2_SCOPES,
		RedirectURL:  config.ROOT_URL + "/oauth2/callback",
		Endpoint: oauth2.Endpoint{
			TokenURL: config.OAUTH2_ENDPOINT_TOKEN_URL,
			AuthURL:  config.OAUTH2_ENDPOINT_AUTH_URL,
		},
	}
}

func (h *BaseHandler) OAuth2LoginGET(w http.ResponseWriter, r *http.Request) {
	conf := newOAuth2(h.config)

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *BaseHandler) OAuth2CallbackGET(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conf := newOAuth2(h.config)

	// Exchange the code for a new token.
	tok, err := conf.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ge the user information.
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(tok))
	resp, err := client.Get(h.config.OAUTH2_ENDPOINT_USER_INFO_URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
