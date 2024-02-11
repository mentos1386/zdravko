package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"code.tjo.space/mentos1386/zdravko/internal"
	"golang.org/x/oauth2"
)

type UserInfo struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

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

func (h *BaseHandler) AuthenticatedUserToOAuth2Token(user *AuthenticatedUser) *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  user.OAuth2AccessToken,
		TokenType:    user.OAuth2TokenType,
		RefreshToken: user.OAuth2RefreshToken,
		Expiry:       user.OAuth2Expiry,
	}
}

func (h *BaseHandler) RefreshToken(w http.ResponseWriter, r *http.Request, user *AuthenticatedUser) (*AuthenticatedUser, error) {
	tok := h.AuthenticatedUserToOAuth2Token(user)
	conf := newOAuth2(h.config)
	refreshed, err := conf.TokenSource(context.Background(), tok).Token()
	if err != nil {
		return nil, err
	}

	refreshedUser := &AuthenticatedUser{
		ID:                 user.ID,
		Email:              user.Email,
		OAuth2AccessToken:  refreshed.AccessToken,
		OAuth2RefreshToken: refreshed.RefreshToken,
		OAuth2TokenType:    refreshed.TokenType,
		OAuth2Expiry:       refreshed.Expiry,
	}

	err = h.SetAuthenticatedUserForRequest(w, r, refreshedUser)
	if err != nil {
		return nil, err
	}

	return refreshedUser, nil
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

	var userInfo UserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.SetAuthenticatedUserForRequest(w, r, &AuthenticatedUser{
		ID:                 userInfo.Sub,
		Email:              userInfo.Email,
		OAuth2AccessToken:  tok.AccessToken,
		OAuth2RefreshToken: tok.RefreshToken,
		OAuth2TokenType:    tok.TokenType,
		OAuth2Expiry:       tok.Expiry,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/settings", http.StatusTemporaryRedirect)
}

func (h *BaseHandler) OAuth2LogoutGET(w http.ResponseWriter, r *http.Request, user *AuthenticatedUser) {
	tok := h.AuthenticatedUserToOAuth2Token(user)
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(tok))
	_, err := client.Get(h.config.OAUTH2_ENDPOINT_USER_INFO_URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.ClearAuthenticatedUserForRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
