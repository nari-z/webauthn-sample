package service

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/duo-labs/webauthn.io/session"
	"github.com/duo-labs/webauthn/protocol"

	"github.com/nari-z/webauthn-sample/domain/repository"
)

type UserRegister struct {
	db repository.UserRepository
	webAuthn *webauthn.WebAuthn
	sessionStore *session.Store
}

func NewUserRegister(db repository.UserRepository, webAuthn *webauthn.WebAuthn, sessionStore *session.Store) *UserRegister {
	return &UserRegister{db: db, webAuthn: webAuthn, sessionStore: sessionStore}
}

func (register *UserRegister) BeginRegistration(c echo.Context) error {
	// get username
	userName := c.Param("username")

	// get user
	user, err := register.db.Find(userName)
	if err != nil {
		displayName := userName
		user = register.db.Create(userName, displayName)
	}

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = user.CredentialExcludeList()
	}

	// generate PublicKeyCredentialCreationOptions, session data
	options, sessionData, err := register.webAuthn.BeginRegistration(
		user,
		registerOptions,
	)
	if err != nil {
		return err
	}

	// store session data as marshaled JSON
	err = register.sessionStore.SaveWebauthnSession("registration", sessionData, c.Request(), c.Response())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, options)
}

func (register *UserRegister) FinishRegistration(c echo.Context) error {
	// get username
	userName := c.Param("username")

	// get user
	user, err := register.db.Find(userName)
	if err != nil {
		return errors.New("not found user")
	}

	// load the session data
	sessionData, err := register.sessionStore.GetWebauthnSession("registration", c.Request())
	if err != nil {
		return err
	}

	credential, err := register.webAuthn.FinishRegistration(user, sessionData, c.Request())
	if err != nil {
		return err
	}

	user.AddCredential(*credential)

	return c.JSON(http.StatusOK, "Registration Success")
}
