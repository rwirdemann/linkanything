package core

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{userRepository: userRepository}

}

func (s UserService) Create(user User) (User, error) {
	return s.userRepository.Create(user)
}

func (s UserService) ByName(name string) (User, error) {
	return s.userRepository.ByName(name)
}

func (s UserService) GetHash(user string) (string, error) {
	return s.userRepository.GetHash(user)
}
