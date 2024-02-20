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

	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

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
	conf := newOAuth2(h.config)

	state := newRandomState()
	result := h.db.Create(&models.OAuth2State{State: state, Expiry: time.Now().Add(5 * time.Minute)})
	if result.Error != nil {
		return result.Error
	}

	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *BaseHandler) OAuth2CallbackGET(c echo.Context) error {
	ctx := context.Background()
	conf := newOAuth2(h.config)

	state := c.QueryParam("state")
	code := c.QueryParam("code")

	result, err := h.query.OAuth2State.WithContext(ctx).Where(
		h.query.OAuth2State.State.Eq(state),
		h.query.OAuth2State.Expiry.Gt(time.Now()),
	).Delete()
	if err != nil {
		return err
	}
	if result.RowsAffected != 1 {
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

	return c.Redirect(http.StatusTemporaryRedirect, "/settings")
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
