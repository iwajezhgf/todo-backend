package storage

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(host, database, user, password string) (*Storage, error) {
	db, err := sql.Open("mysql", user+":"+password+"@("+host+")/"+database+"?parseTime=true")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	mysql := &Storage{DB: db}

	return mysql, nil
}
