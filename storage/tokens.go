package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/iwajezhgf/todo-backend/types"
)

type TokenStorage struct {
	*Storage
}

func (s *TokenStorage) Create(token string, userId uint64) error {
	_, err := s.DB.Exec(
		"INSERT INTO tokens (token, expire, user_id) VALUES (?, ?, ?)", token,
		time.Now().UTC().Add(30*24*time.Hour), userId,
	)
	if err != nil {
		log.Println("mysql error: ", err)
		return errors.New("mysql error")
	}

	return nil
}

func (s *TokenStorage) GetByToken(token string) (*types.Token, error) {
	var t types.Token
	err := s.DB.QueryRow("SELECT id, token, user_id FROM tokens WHERE token = ?", token).Scan(
		&t.ID, &t.Token, &t.UserID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return &t, nil
}

func (s *TokenStorage) Delete(token string) {
	_, err := s.DB.Exec("DELETE FROM tokens WHERE token = ?", token)
	if err != nil {
		log.Println("mysql error: ", err)
	}
}

func (s *TokenStorage) StartTokenCleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.CleanupExpiredTokens()
		}
	}
}

func (s *TokenStorage) CleanupExpiredTokens() {
	_, err := s.DB.Exec("DELETE FROM tokens WHERE expire < ?", time.Now().UTC())
	if err != nil {
		log.Println("mysql error: ", err)
	}
}
