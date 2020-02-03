package repository

import (
	"github.com/nari-z/webauthn-sample/domain/model"
)

type UserRepository interface {
	Register(user model.User) error
	Find(userName string) (*model.User, error)
	Create(userName string, displayName string) *model.User
	Update(model.User) error
}