package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mentos1386/zdravko/database/models"
	"github.com/mentos1386/zdravko/internal/config"
	"github.com/mentos1386/zdravko/internal/server/services"
	"golang.org/x/oauth2"
)

const oauth2RedirectSessionName = "zdravko-hey-oauth2"

func (h *BaseHandler) setOAuth2Redirect(c echo.Context, redirect string) error {
	w := c.Response()
	r := c.Request()

	session, err := h.store.Get(r, oauth2RedirectSessionName)
	if err != nil {
		return err
	}
	session.Values["redirect"] = redirect
	return h.store.Save(r, w, session)
}

func (h *BaseHandler) getOAuth2Redirect(c echo.Context) (string, error) {
	r := c.Request()

	session, err := h.store.Get(r, oauth2RedirectSessionName)
	if err != nil {
		return "", err
	}
	if session.IsNew {
		return "", nil
	}
	return session.Values["redirect"].(string), nil
}

func (h *BaseHandler) clearOAuth2Redirect(c echo.Context) error {
	w := c.Response()
	r := c.Request()

	session, err := h.store.Get(r, oauth2RedirectSessionName)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return h.store.Save(r, w, session)
}

type UserInfo struct {
	Id    int    `json:"id"` // FIXME: This might not always be int?
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

func newRandomState() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func newOAuth2(config *config.ServerConfig) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.OAuth2.ClientID,
		ClientSecret: config.OAuth2.ClientSecret,
		Scopes:       config.OAuth2.Scopes,
		RedirectURL:  config.RootUrl + "/oauth2/callback",
		Endpoint: oauth2.Endpoint{
			TokenURL: config.OAuth2.EndpointTokenURL,
			AuthURL:  config.OAuth2.EndpointAuthURL,
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
		fmt.Println("Error: ", err)
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

func (h *BaseHandler) OAuth2LoginGET(c echo.Context) error {
	ctx := context.Background()
	conf := newOAuth2(h.config)

	state := newRandomState()
	err := services.CreateOAuth2State(ctx, h.db, &models.OAuth2State{
		State:     state,
		ExpiresAt: &models.Time{Time: time.Now().Add(5 * time.Minute)},
	})
	if err != nil {
		return err
	}

	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)

	redirect := c.QueryParam("redirect")
	h.logger.Info("OAuth2LoginGET", "redirect", redirect)

	err = h.setOAuth2Redirect(c, redirect)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *BaseHandler) OAuth2CallbackGET(c echo.Context) error {
	ctx := context.Background()
	conf := newOAuth2(h.config)

	state := c.QueryParam("state")
	code := c.QueryParam("code")

	deleted, err := services.DeleteOAuth2State(ctx, h.db, state)
	if err != nil {
		return err
	}
	if !deleted {
		return errors.New("invalid state")
	}

	// Exchange the code for a new token.
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		return err
	}

	// Ge the user information.
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(tok))
	resp, err := client.Get(h.config.OAuth2.EndpointUserInfoURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var userInfo UserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return err
	}

	userId := userInfo.Sub
	if userInfo.Id != 0 {
		userId = strconv.Itoa(userInfo.Id)
	}

	err = h.SetAuthenticatedUserForRequest(c.Response(), c.Request(), &AuthenticatedUser{
		ID:                 userId,
		Email:              userInfo.Email,
		OAuth2AccessToken:  tok.AccessToken,
		OAuth2RefreshToken: tok.RefreshToken,
		OAuth2TokenType:    tok.TokenType,
		OAuth2Expiry:       tok.Expiry,
	})
	if err != nil {
		return err
	}

	redirect, err := h.getOAuth2Redirect(c)
	if err != nil {
		return err
	}
	h.logger.Info("OAuth2CallbackGET", "redirect", redirect)
	if redirect == "" {
		redirect = "/settings"
	}

	err = h.clearOAuth2Redirect(c)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, redirect)
}

func (h *BaseHandler) OAuth2LogoutGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	if h.config.OAuth2.EndpointLogoutURL != "" {
		tok := h.AuthenticatedUserToOAuth2Token(cc.Principal.User)
		client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(tok))
		_, err := client.Get(h.config.OAuth2.EndpointLogoutURL)
		if err != nil {
			return err
		}
	}

	err := h.ClearAuthenticatedUserForRequest(c.Response(), c.Request())
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
