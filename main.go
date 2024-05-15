package main

import (
	"fmt"
	"log"

	"github.com/iwajezhgf/todo-backend/api"
	"github.com/iwajezhgf/todo-backend/storage"
	"github.com/spf13/viper"
)

type config struct {
	DB     dbConfig     `yaml:"db"`
	Server serverConfig `yaml:"server"`
}

type dbConfig struct {
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type serverConfig struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

func main() {
	var c config
	if err := initConfig("config.yml", &c); err != nil {
		log.Fatalf("failed to load config file: %s", err)
	}

	db, err := storage.NewStorage(c.DB.Host, c.DB.Database, c.DB.User, c.DB.Password)
	if err != nil {
		log.Fatalf("db error: %s", err)
	}

	todos := &storage.TodoStorage{Storage: db}
	tokens := &storage.TokenStorage{Storage: db}
	users := &storage.UserStorage{Storage: db}

	go todos.StartTodoStatus()
	go tokens.StartTokenCleanup()

	serverHost := fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)

	s := api.Server{
		Todos:  todos,
		Tokens: tokens,
		Users:  users,
	}
	if err = s.Run(serverHost); err != nil {
		log.Fatalf("fasthttp server error: %s", err)
	}
}

func initConfig(path string, conf *config) error {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(conf); err != nil {
		return err
	}

	return nil
}
