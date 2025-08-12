package repositories

import (
	"fmt"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetUserByName(name string) string {
	query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)

	// запрос в базу
	fmt.Println(query)

	return "test name"
}
