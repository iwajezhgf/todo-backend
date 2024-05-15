package types

import "time"

type User struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
}

type Token struct {
	ID     uint64    `json:"id"`
	Token  string    `json:"token"`
	Expire time.Time `json:"expire"`
	UserID uint64    `json:"user_id"`
}

type Todo struct {
	ID      uint64    `json:"id"`
	Title   string    `json:"title"`
	Note    string    `json:"note"`
	Created time.Time `json:"created"`
	Expire  time.Time `json:"expire"`
	Status  string    `json:"status"`
	UserID  uint64    `json:"user_id"`
}

type TodoData struct {
	Items      []Todo     `json:"items"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Prev []int `json:"prev"`
	Next []int `json:"next"`
}
