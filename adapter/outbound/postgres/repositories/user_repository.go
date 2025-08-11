package repositories

import (
	"database/sql"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	// TODO: В реальном приложении здесь должна быть инициализация подключения к БД
	return &UserRepository{
		db: nil, // Пока используем nil, т.к. это mock
	}
}

func (r *UserRepository) GetUserByName(name string) string {
	// БЕЗОПАСНЫЙ ВАРИАНТ - используем параметризованный запрос
	// query := "SELECT * FROM users WHERE name = $1"
	// rows, err := r.db.Query(query, name)
	
	// Текущий mock-реализация для демонстрации
	fmt.Printf("Безопасный запрос: SELECT * FROM users WHERE name = $1 (параметр: %s)\n", name)
	
	return "test name"
}
