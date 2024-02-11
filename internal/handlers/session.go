package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const sessionName = "zdravko-hey"

type AuthenticatedUser struct {
	ID                 string
	Email              string
	OAuth2AccessToken  string
	OAuth2RefreshToken string
	OAuth2TokenType    string
	OAuth2Expiry       time.Time
}

type authenticatedUserKeyType string

const authenticatedUserKey authenticatedUserKeyType = "authenticatedUser"

func WithUser(ctx context.Context, user *AuthenticatedUser) context.Context {
	return context.WithValue(ctx, authenticatedUserKey, user)
}

func GetUser(ctx context.Context) *AuthenticatedUser {
	user, ok := ctx.Value(authenticatedUserKey).(*AuthenticatedUser)
	if !ok {
		return nil
	}
	return user
}

func (h *BaseHandler) GetAuthenticatedUserForRequest(r *http.Request) (*AuthenticatedUser, error) {
	session, err := h.store.Get(r, sessionName)
	if err != nil {
		return nil, err
	}
	if session.IsNew {
		return nil, fmt.Errorf("session is nil")
	}

	expiry, err := time.Parse(time.RFC3339, session.Values["oauth2_expiry"].(string))
	if err != nil {
		return nil, err
	}

	user := &AuthenticatedUser{
		ID:                 session.Values["id"].(string),
		Email:              session.Values["email"].(string),
		OAuth2AccessToken:  session.Values["oauth2_access_token"].(string),
		OAuth2RefreshToken: session.Values["oauth2_refresh_token"].(string),
		OAuth2TokenType:    session.Values["oauth2_token_type"].(string),
		OAuth2Expiry:       expiry,
	}

	return user, nil
}

func (h *BaseHandler) SetAuthenticatedUserForRequest(w http.ResponseWriter, r *http.Request, user *AuthenticatedUser) error {
	session, err := h.store.Get(r, sessionName)
	if err != nil {
		return err
	}
	session.Values["id"] = user.ID
	session.Values["email"] = user.Email
	session.Values["oauth2_access_token"] = user.OAuth2AccessToken
	session.Values["oauth2_refresh_token"] = user.OAuth2RefreshToken
	session.Values["oauth2_token_type"] = user.OAuth2TokenType
	session.Values["oauth2_expiry"] = user.OAuth2Expiry.Format(time.RFC3339)
	err = h.store.Save(r, w, session)
	if err != nil {
		return err
	}
	return nil
}

func (h *BaseHandler) ClearAuthenticatedUserForRequest(w http.ResponseWriter, r *http.Request) error {
	session, err := h.store.Get(r, sessionName)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	err = h.store.Save(r, w, session)
	if err != nil {
		return err
	}
	return nil
}

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *AuthenticatedUser)

func (h *BaseHandler) Authenticated(next AuthenticatedHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := h.GetAuthenticatedUserForRequest(r)
		if err != nil {
			http.Redirect(w, r, "/oauth2/login", http.StatusTemporaryRedirect)
			return
		}
		if user.OAuth2Expiry.Before(time.Now()) {
			user, err = h.RefreshToken(w, r, user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		next(w, r, user)
	}
}
