package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/iwajezhgf/todo-backend/types"
	"github.com/valyala/fasthttp"
)

const (
	AlreadyExists      = "already_exists"
	InvalidCredentials = "invalid_credentials"
	InvalidEmail       = "invalid_email"
	InvalidPassword    = "invalid_password"
	InvalidDate        = "invalid_date"
	Unauthorized       = "unauthorized"
	NotFound           = "not_found"
	PasswordEquals     = "password_equals"
)

func okResponse(ctx *fasthttp.RequestCtx, code int, response interface{}) {
	jsonResponse := map[string]interface{}{
		"ok":       true,
		"response": response,
	}
	jsonResponseBytes, err := json.Marshal(jsonResponse)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusBadRequest, "", "")
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(code)
	ctx.SetBody(jsonResponseBytes)
}

func errorResponse(ctx *fasthttp.RequestCtx, statusCode int, errorMessage string, errorCode string) {
	response := map[string]interface{}{
		"error": true,
		"response": map[string]interface{}{
			"code":    errorCode,
			"message": errorMessage,
			"path":    string(ctx.Path()),
		},
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		ctx.SetContentType("text/plain")
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte("Internal Server Error"))
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)
	ctx.SetBody(responseBytes)
}

func (s *Server) getUserByToken(ctx *fasthttp.RequestCtx) (*types.User, error) {
	token := ctx.UserValue("token").(string)
	t, err := s.Tokens.GetByToken(token)
	if err != nil {
		return nil, err
	}

	return s.Users.GetById(t.UserID)
}

func (s *Server) validateFields(entity interface{}) func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			var data map[string]interface{}
			err := json.Unmarshal(ctx.PostBody(), &data)
			if err != nil {
				errorResponse(ctx, fasthttp.StatusBadRequest, "Invalid body", "")
				return
			}

			entityType := reflect.TypeOf(entity)
			if entityType.Kind() != reflect.Struct {
				errorResponse(ctx, fasthttp.StatusInternalServerError, "Invalid entity type", "")
				return
			}

			for i := 0; i < entityType.NumField(); i++ {
				field := entityType.Field(i)
				fieldName := strings.ToLower(field.Tag.Get("json"))
				if value, ok := data[fieldName]; !ok || fmt.Sprintf("%v", value) == "" {
					errorResponse(ctx, fasthttp.StatusBadRequest, fmt.Sprintf("Invalid %s", fieldName), "")
					return
				}
			}

			next(ctx)
		}
	}
}

func (s *Server) onlyAuthorized(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		authHeader := ctx.Request.Header.Peek("Authorization")
		if authHeader == nil {
			errorResponse(ctx, fasthttp.StatusUnauthorized, "Unauthorized", Unauthorized)
			return
		}

		const bearer = "Bearer "
		if !bytes.HasPrefix(authHeader, []byte(bearer)) {
			errorResponse(ctx, fasthttp.StatusUnauthorized, "Unauthorized", Unauthorized)
			return
		}

		token := string(authHeader[len(bearer):])
		ctx.SetUserValue("token", token)
		next(ctx)
	}
}
