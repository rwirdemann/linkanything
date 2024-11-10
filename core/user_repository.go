package core

type UserRepository interface {
	Create(user User) (User, error)
	ByName(name string) (User, error)
	GetHash(user string) (string, error)
	DeleteAll() error
}
