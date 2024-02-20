package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwtInternal "code.tjo.space/mentos1386/zdravko/internal/jwt"
	"github.com/labstack/echo/v4"
)

const sessionName = "zdravko-hey"

type AuthenticatedPrincipal struct {
	User   *AuthenticatedUser
	Worker *AuthenticatedWorker
}

type AuthenticatedUser struct {
	ID                 string
	Email              string
	OAuth2AccessToken  string
	OAuth2RefreshToken string
	OAuth2TokenType    string
	OAuth2Expiry       time.Time
}

type AuthenticatedWorker struct {
	Slug  string
	Group string
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

func (h *BaseHandler) AuthenticateRequestWithCookies(r *http.Request) (*AuthenticatedUser, error) {
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

func (h *BaseHandler) AuthenticateRequestWithToken(r *http.Request) (*AuthenticatedPrincipal, error) {
	authorization := r.Header.Get("Authorization")

	splitAuthorization := strings.Split(authorization, " ")
	if len(splitAuthorization) != 2 {
		return nil, fmt.Errorf("invalid authorization header")
	}
	if splitAuthorization[0] != "Bearer" {
		return nil, fmt.Errorf("invalid authorization header")
	}

	_, claims, err := jwtInternal.ParseToken(splitAuthorization[1], h.config.Jwt.PublicKey)
	if err != nil {
		return nil, err
	}

	splitSubject := strings.Split(claims.Subject, ":")
	if len(splitSubject) != 2 {
		return nil, fmt.Errorf("invalid subject")
	}

	var worker *AuthenticatedWorker
	var user *AuthenticatedUser

	if splitSubject[0] == "user" {
		user = &AuthenticatedUser{}
	} else if splitSubject[0] == "worker" {
		worker = &AuthenticatedWorker{
			Slug:  splitSubject[1],
			Group: claims.WorkerGroup,
		}
	}

	principal := &AuthenticatedPrincipal{
		User:   user,
		Worker: worker,
	}

	return principal, nil
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

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *AuthenticatedPrincipal)

type AuthenticatedContext struct {
	echo.Context
	Principal *AuthenticatedPrincipal
}

func (h *BaseHandler) Authenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// First try cookie authentication
		user, err := h.AuthenticateRequestWithCookies(c.Request())
		if err == nil {
			if user.OAuth2Expiry.Before(time.Now()) {
				user, err = h.RefreshToken(c.Response(), c.Request(), user)
				if err != nil {
					return c.Redirect(http.StatusTemporaryRedirect, "/oauth2/login")
				}
			}

			cc := AuthenticatedContext{c, &AuthenticatedPrincipal{user, nil}}
			return next(cc)
		}
		// Then try token based authentication
		principal, err := h.AuthenticateRequestWithToken(c.Request())
		if err == nil {
			cc := AuthenticatedContext{c, principal}
			return next(cc)
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/oauth2/login")
	}
}
