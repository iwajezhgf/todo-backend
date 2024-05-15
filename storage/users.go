package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/iwajezhgf/todo-backend/types"
)

type UserStorage struct {
	*Storage
}

func (s *UserStorage) Create(email string, password []byte) error {
	_, err := s.DB.Exec("INSERT INTO users (email, password) VALUES (?, ?)", email, password)
	if err != nil {
		log.Println("mysql error: ", err)
		return errors.New("mysql error")
	}

	return nil
}

func (s *UserStorage) GetById(id uint64) (*types.User, error) {
	var user types.User
	err := s.DB.QueryRow("SELECT id, email, password FROM users WHERE id = ?", id).Scan(
		&user.ID, &user.Email, &user.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (s *UserStorage) GetByEmail(email string) (*types.User, error) {
	var user types.User
	err := s.DB.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).Scan(
		&user.ID, &user.Email, &user.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (s *UserStorage) EditPassword(userId uint64, password []byte) error {
	_, err := s.DB.Exec(
		"UPDATE users SET password = ? WHERE id = ?",
		password,
		userId,
	)
	if err != nil {
		log.Println("mysql error: ", err)
		return errors.New("mysql error")
	}

	return nil
}
