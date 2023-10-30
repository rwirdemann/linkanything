package service

import (
	"github.com/rwirdemann/linkanything/core/domain"
	"github.com/rwirdemann/linkanything/core/port"
)

type UserService struct {
	userRepository port.UserRepository
}

func NewUserService(userRepository port.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}

}

func (s UserService) Create(user domain.User) (domain.User, error) {
	return s.userRepository.Create(user)
}
