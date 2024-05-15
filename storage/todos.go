package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/iwajezhgf/todo-backend/types"
)

type TodoStorage struct {
	*Storage
}

func (s *TodoStorage) Create(title, note string, expire time.Time, userId uint64) error {
	_, err := s.DB.Exec(
		"INSERT INTO todos (title, note, created, expire, user_id) VALUES (?, ?, ?, ?, ?)", title, note,
		time.Now().UTC(), expire, userId,
	)
	if err != nil {
		log.Println("mysql error: ", err)
		return errors.New("mysql error")
	}

	return nil
}

func (s *TodoStorage) GetByUserId(id, userId uint64) (*types.Todo, error) {
	var t types.Todo
	err := s.DB.QueryRow("SELECT id, status, expire, user_id FROM todos WHERE user_id = ? AND id = ?", userId, id).Scan(
		&t.ID, &t.Status, &t.Expire, &t.UserID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &t, nil
}

func (s *TodoStorage) GetTodosByPage(userId uint64, page, limit int) (*types.TodoData, error) {
	queryCount := "SELECT COUNT(*) FROM todos WHERE user_id = ?"
	var totalRecords int
	err := s.DB.QueryRow(queryCount, userId).Scan(&totalRecords)
	if err != nil {
		log.Println("mysql error: ", err)
		return nil, errors.New("mysql error")
	}

	totalPages := (totalRecords + limit - 1) / limit

	offset := (page - 1) * limit
	query := "SELECT * FROM todos WHERE user_id = ? ORDER BY created DESC LIMIT ? OFFSET ?"

	rows, err := s.DB.Query(query, userId, limit, offset)
	if err != nil {
		log.Println("mysql error: ", err)
		return nil, errors.New("mysql error")
	}
	defer rows.Close()

	todos := make([]types.Todo, 0)
	for rows.Next() {
		var todo types.Todo
		err = rows.Scan(&todo.ID, &todo.Title, &todo.Note, &todo.Created, &todo.Expire, &todo.Status, &todo.UserID)
		if err != nil {
			log.Println("mysql error: ", err)
			continue
		}
		todos = append(todos, todo)
	}

	var data types.TodoData
	data.Items = todos

	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 0
	}

	prevPages := []int{}
	if prevPage > 0 {
		prevPages = append(prevPages, prevPage)
	}

	nextPage := page + 1

	nextPages := []int{}
	if nextPage <= totalPages {
		nextPages = append(nextPages, nextPage)
	}

	data.Pagination.Prev = prevPages
	data.Pagination.Next = nextPages

	return &data, nil
}

func (s *TodoStorage) Edit(title, note string, expire time.Time, id, userId uint64) error {
	_, err := s.DB.Exec(
		"UPDATE todos SET title = ?, note = ?, expire = ? WHERE id = ? AND user_id = ?",
		title,
		note,
		expire,
		id,
		userId,
	)
	if err != nil {
		log.Println("mysql error: ", err)
		return errors.New("mysql error")
	}

	return nil
}

func (s *TodoStorage) EditStatus(id, userId uint64, status string) error {
	_, err := s.DB.Exec(
		"UPDATE todos SET status = ? WHERE id = ? AND user_id = ?",
		status,
		id,
		userId,
	)
	if err != nil {
		log.Println("mysql error: ", err)
		return errors.New("mysql error")
	}

	return nil
}

func (s *TodoStorage) Delete(id uint64) {
	_, err := s.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		log.Println("mysql error: ", err)
	}
}

func (s *TodoStorage) StartTodoStatus() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.CheckTodoStatus()
		}
	}
}

func (s *TodoStorage) CheckTodoStatus() {
	rows, err := s.DB.Query("SELECT id, expire FROM todos WHERE status != 'completed'")
	if err != nil {
		log.Println("mysql error: ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id uint64
		var expire time.Time
		if err = rows.Scan(&id, &expire); err != nil {
			log.Println("mysql error: ", err)
			continue
		}

		if time.Now().After(expire) {
			if _, err = s.DB.Exec("UPDATE todos SET status = 'overdue' WHERE id = ?", id); err != nil {
				log.Println("mysql error: ", err)
				continue
			}
		}
	}
	if err = rows.Err(); err != nil {
		log.Println("mysql error: ", err)
	}
}
