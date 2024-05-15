package api

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type createTodo struct {
	Title  string `json:"title"`
	Note   string `json:"note"`
	Expire string `json:"expire"`
}

func (s *Server) handleCreateTodo(ctx *fasthttp.RequestCtx) {
	var todo createTodo
	json.Unmarshal(ctx.PostBody(), &todo)

	expireTime, err := time.Parse("2006-01-02 15:04:05", todo.Expire)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusBadRequest, "", InvalidDate)
		return
	}

	u, err := s.getUserByToken(ctx)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusUnauthorized, "Unauthorized", Unauthorized)
		return
	}

	err = s.Todos.Create(todo.Title, todo.Note, expireTime, u.ID)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusInternalServerError, err.Error(), "")
		return
	}

	okResponse(ctx, fasthttp.StatusOK, map[string]string{})
}

func (s *Server) handleGetTodos(ctx *fasthttp.RequestCtx) {
	queryParams := string(ctx.QueryArgs().QueryString())
	if queryParams == "" {
		errorResponse(ctx, fasthttp.StatusBadRequest, "Invalid query params", "")
		return
	}

	queryParts := strings.Split(queryParams, "&")
	if len(queryParts) != 2 {
		errorResponse(ctx, fasthttp.StatusBadRequest, "Invalid query params", "")
		return
	}

	var limit, page int
	for _, part := range queryParts {
		keyValue := strings.Split(part, "=")
		if len(keyValue) != 2 {
			continue
		}
		key := keyValue[0]
		value := keyValue[1]

		if key == "limit" {
			limit, _ = strconv.Atoi(value)
		} else if key == "page" {
			page, _ = strconv.Atoi(value)
		}
	}

	if limit == 0 {
		limit = 10
	}
	if page == 0 {
		page = 1
	}

	u, err := s.getUserByToken(ctx)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusUnauthorized, "Unauthorized", Unauthorized)
		return
	}

	todoResponse, err := s.Todos.GetTodosByPage(u.ID, page, limit)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusInternalServerError, err.Error(), "")
		return
	}

	okResponse(ctx, fasthttp.StatusOK, todoResponse)
}

type editTodo struct {
	ID     uint64 `json:"id"`
	Title  string `json:"title"`
	Note   string `json:"note"`
	Expire string `json:"expire"`
}

func (s *Server) handleEditTodo(ctx *fasthttp.RequestCtx) {
	var todo editTodo
	json.Unmarshal(ctx.PostBody(), &todo)

	expireTime, err := time.Parse("2006-01-02 15:04:05", todo.Expire)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusBadRequest, "", InvalidDate)
		return
	}

	u, err := s.getUserByToken(ctx)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusUnauthorized, "Unauthorized", Unauthorized)
		return
	}

	_, err = s.Todos.GetByUserId(todo.ID, u.ID)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusNotFound, "", NotFound)
		return
	}

	err = s.Todos.Edit(todo.Title, todo.Note, expireTime, todo.ID, u.ID)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusInternalServerError, err.Error(), "")
		return
	}

	okResponse(ctx, fasthttp.StatusOK, map[string]string{})
}

func (s *Server) handleCompleteTodo(ctx *fasthttp.RequestCtx) {
	u, err := s.getUserByToken(ctx)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusUnauthorized, "Unauthorized", Unauthorized)
		return
	}

	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusBadRequest, "Invalid ID", "")
		return
	}

	todo, err := s.Todos.GetByUserId(id, u.ID)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusNotFound, "", NotFound)
		return
	}

	if todo.Status != "completed" {
		todo.Status = "completed"
	} else {
		if time.Now().After(todo.Expire) {
			todo.Status = "overdue"
		} else {
			todo.Status = "active"
		}
	}

	err = s.Todos.EditStatus(id, u.ID, todo.Status)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusInternalServerError, err.Error(), "")
		return
	}

	okResponse(ctx, fasthttp.StatusOK, map[string]string{})
}

func (s *Server) handleDeleteTodo(ctx *fasthttp.RequestCtx) {
	u, err := s.getUserByToken(ctx)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusUnauthorized, "Unauthorized", Unauthorized)
		return
	}

	idStr := ctx.UserValue("id").(string)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusBadRequest, "Invalid ID", "")
		return
	}

	_, err = s.Todos.GetByUserId(id, u.ID)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusNotFound, "", NotFound)
		return
	}

	s.Todos.Delete(id)
	okResponse(ctx, fasthttp.StatusOK, map[string]string{})
}
