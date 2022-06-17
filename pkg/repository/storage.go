package repository

import (
	"errors"
	"jwt-auth-gin/pkg/api"
)

type Storage interface {
	GetUser(userName string) (user api.User, err error)
}

type storage struct {
	jwtKey []byte
	users  map[string]string
}

func NewStorage(jwtKey []byte, users map[string]string) Storage {
	return &storage{
		jwtKey: jwtKey,
		users:  users,
	}
}

func (s *storage) GetUser(username string) (user api.User, err error) {
	pw, exists := s.users[username]

	if !exists {
		err = errors.New("User does not exist")
		return api.User{}, err
	}

	user.Username = username
	user.Password = pw
	return user, nil
}
