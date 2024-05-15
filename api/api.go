package api

import (
	"time"

	"github.com/fasthttp/router"
	"github.com/iwajezhgf/todo-backend/storage"
	"github.com/valyala/fasthttp"
)

type Server struct {
	s *fasthttp.Server

	Todos  *storage.TodoStorage
	Tokens *storage.TokenStorage
	Users  *storage.UserStorage
}

func (s *Server) Run(host string) error {
	r := router.New()

	api := r.Group("/api")
	{
		api.POST("/register", s.validateFields(authUser{})(s.handleRegister))
		api.POST("/login", s.validateFields(authUser{})(s.handleLogin))
		api.GET("/auth", s.onlyAuthorized(s.handleAuth))
		api.POST("/logout", s.onlyAuthorized(s.handleLogout))

		api.POST("/settings/password", s.onlyAuthorized(s.validateFields(settingsPassword{})(s.ChangePassword)))

		api.POST("/todo", s.onlyAuthorized(s.validateFields(createTodo{})(s.handleCreateTodo)))
		api.POST("/todo/{id}", s.onlyAuthorized(s.handleCompleteTodo))
		api.GET("/todo", s.onlyAuthorized(s.handleGetTodos))
		api.PUT("/todo", s.onlyAuthorized(s.validateFields(editTodo{})(s.handleEditTodo)))
		api.DELETE("/todo/{id}", s.onlyAuthorized(s.handleDeleteTodo))
	}

	srv := &fasthttp.Server{
		Handler:           r.Handler,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		ReduceMemoryUsage: true,
	}

	return srv.ListenAndServe(host)
}
