package mock

import (
	"encoding/binary"
	"errors"
	"github.com/nari-z/webauthn-sample/domain/model"
	"github.com/nari-z/webauthn-sample/domain/repository"
)

type UserRepository struct {
	users map[string]*model.User
	userCount uint64
}

func NewUserRepository() repository.UserRepository {
	return &UserRepository{users: make(map[string]*model.User), userCount: 0}
}

func (repository *UserRepository) Register(user model.User) error {
	repository.users[string(user.WebAuthnID())] = &user
	return nil
}

func (repository UserRepository) Find(userName string) (*model.User, error) {
	for _, user := range repository.users {
		if user.WebAuthnName() == userName {
			return user, nil
		}
	}

	return nil, errors.New("not found user")
}

func (repository UserRepository) Create(userName string, displayName string) *model.User {
	newUser := model.NewUser(repository.createID(), userName, displayName, "")
	repository.Register(*newUser)
	return newUser
}

func (repository *UserRepository) Update(user model.User) error {
	_, ok := repository.users[string(user.WebAuthnID())]
	if ok == false {
		return errors.New("not found user")
	}

	return nil
}

func (repository UserRepository) createID() []byte {
    buf := make([]byte, binary.MaxVarintLen64) //MaxVarintLen64 = 10
	binary.BigEndian.PutUint64(buf, repository.userCount)
	repository.userCount = repository.userCount + 1
	return buf
}
