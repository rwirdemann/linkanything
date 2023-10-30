package port

import "github.com/rwirdemann/linkanything/core/domain"

type UserService interface {
	Create(user domain.User) (domain.User, error)
}
