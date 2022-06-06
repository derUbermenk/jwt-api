package api

import "log"

// UserService contains the methods of the user service
// declaring this interface allows us to mock this interface
// for testing
type UserService interface {
	GetUser(userName string) (user User, err error)
}

type userRepository interface {
	GetUser(username string) (user User, err error)
}

type userService struct {
	storage userRepository
}

func NewUserService(userRepo userRepository) UserService {
	return &userService{
		storage: userRepo,
	}
}

func (u *userService) GetUser(username string) (user User, err error) {
	user, err = u.storage.GetUser(username)

	if err != nil {
		log.Printf("Service error: %v", err)
		return user, err
	}

	return user, nil
}
