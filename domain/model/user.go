package model

import (
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/duo-labs/webauthn/protocol"
)

type User struct {
	id []byte
	name string
	displayName string
	icon string
	credentials []webauthn.Credential
}

func NewUser(id []byte, name string, displayName string, icon string) *User {
	return &User{id: id, name: name, displayName: displayName, icon: icon}
}

func (user User) WebAuthnID() []byte {
	return user.id
}

func (user User) WebAuthnName() string {
	return user.name
}

func (user User) WebAuthnDisplayName() string {
	return user.displayName
}

func (user User) WebAuthnIcon() string {
	return user.icon
}

func (user User) WebAuthnCredentials() []webauthn.Credential {
	return user.credentials
}

func (user *User) AddCredential(cred webauthn.Credential) {
	user.credentials = append(user.credentials, cred)
}

func (user User) CredentialExcludeList() []protocol.CredentialDescriptor {
	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range user.credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}