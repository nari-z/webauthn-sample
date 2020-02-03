package main

import (
	"fmt"
	"log"

	"github.com/duo-labs/webauthn/webauthn"

	"github.com/duo-labs/webauthn.io/session"
	"github.com/nari-z/webauthn-sample/datastore/mock"
	"github.com/nari-z/webauthn-sample/domain/repository"
	"github.com/nari-z/webauthn-sample/provider"
	"github.com/nari-z/webauthn-sample/service"
)

func main() {
	fmt.Println("Hello WebAuthn.")

	// setting webauthn.
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "webauthn-sample",
		RPID:          "localhost",
		RPOrigin:      "http://localhost:4444",
		// RPIcon: "https://duo.com/logo.png",
	})

	// create session store.
	sessionStore, err := session.NewStore()
	if err != nil {
		log.Fatal("failed to create session store:", err)
	}

	// create db.
	var db repository.UserRepository
	db = mock.NewUserRepository()

	// create services.
	userRegister := service.NewUserRegister(db, webAuthn, sessionStore)
	userAuthorizer := service.NewUserAuthorizer(db, webAuthn, sessionStore)

	// launch host.
	p := provider.NewProvider(4444)
	p.UseCORS()
	p.RegisterSamplePage()
	p.RegisterGetMethod("/attenstation/options/:username", userRegister.BeginRegistration)
	p.RegisterPostMethod("/attenstation/result/:username", userRegister.FinishRegistration)
	p.RegisterGetMethod("/assertion/options/:username", userAuthorizer.Begin)
	p.RegisterPostMethod("/assertion/result/:username", userAuthorizer.Finish)
	
	p.Run()
}
