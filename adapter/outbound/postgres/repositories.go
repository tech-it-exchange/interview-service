package postgres

import (
	"interview-service/adapter/outbound/postgres/repositories"
)

type Repositories struct {
	User *repositories.UserRepository
}

func NewRepositories() *Repositories {
	user := repositories.NewUserRepository()

	return &Repositories{
		User: user,
	}
}
