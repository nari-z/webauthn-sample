package service

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/duo-labs/webauthn.io/session"

	"github.com/nari-z/webauthn-sample/domain/repository"
)

type UserAuthorizer struct {
	db repository.UserRepository
	webAuthn *webauthn.WebAuthn
	sessionStore *session.Store
}

func NewUserAuthorizer(db repository.UserRepository, webAuthn *webauthn.WebAuthn, sessionStore *session.Store) *UserAuthorizer {
	return &UserAuthorizer{db: db, webAuthn: webAuthn, sessionStore: sessionStore}
}

func (authorizer *UserAuthorizer) Begin(c echo.Context) error {
	// get username
	userName := c.Param("username")

	// get user
	user, err := authorizer.db.Find(userName)
	if err != nil {
		fmt.Println("user not found")
		return err
	}

	// generate PublicKeyCredentialRequestOptions, session data
	options, sessionData, err := authorizer.webAuthn.BeginLogin(user)
	if err != nil {
		fmt.Println("BeginLogin")
		return err
	}

	// store session data as marshaled JSON
	err = authorizer.sessionStore.SaveWebauthnSession("authentication", sessionData, c.Request(), c.Response())
	if err != nil {
		fmt.Println("SaveWebauthnSession")
		return err
	}

	return c.JSON(http.StatusOK, options)
}

func (authorizer *UserAuthorizer) Finish(c echo.Context) error {
	// get username
	userName := c.Param("username")

	// get user
	user, err := authorizer.db.Find(userName)
	if err != nil {
		return err
	}

	// load the session data
	sessionData, err := authorizer.sessionStore.GetWebauthnSession("authentication", c.Request())
	if err != nil {
		fmt.Println("GetWebauthnSession")
		return err
	}

	_, err = authorizer.webAuthn.FinishLogin(user, sessionData, c.Request())
	if err != nil {
		fmt.Println("FinishLogin")
		return err
	}

	// success
	return c.JSON(http.StatusOK, "OK")
}
