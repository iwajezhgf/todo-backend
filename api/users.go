package api

import (
	"encoding/json"
	"math/rand"
	"regexp"

	"github.com/iwajezhgf/todo-backend/types"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"
)

type authUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) handleRegister(ctx *fasthttp.RequestCtx) {
	var u authUser
	json.Unmarshal(ctx.PostBody(), &u)

	emailRegexp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegexp.MatchString(u.Email) {
		errorResponse(ctx, fasthttp.StatusBadRequest, "", InvalidEmail)
		return
	}

	passRegexp := regexp.MustCompile(`^.{6,}$`)
	if !passRegexp.MatchString(u.Password) {
		errorResponse(ctx, fasthttp.StatusBadRequest, "", InvalidPassword)
		return
	}

	_, err := s.Users.GetByEmail(u.Email)
	if err == nil {
		errorResponse(ctx, fasthttp.StatusBadRequest, "", AlreadyExists)
		return
	}

	hashedPassword, _ := hashPassword(u.Password)

	err = s.Users.Create(u.Email, hashedPassword)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusInternalServerError, err.Error(), "")
		return
	}

	okResponse(ctx, fasthttp.StatusCreated, map[string]string{})
}

type authResponse struct {
	Auth  types.User `json:"auth"`
	Token string     `json:"token"`
}

func (s *Server) handleLogin(ctx *fasthttp.RequestCtx) {
	var u authUser
	json.Unmarshal(ctx.PostBody(), &u)

	user, err := s.Users.GetByEmail(u.Email)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusBadRequest, "", InvalidCredentials)
		return
	}

	if !checkPasswordHash(u.Password, user.Password) {
		errorResponse(ctx, fasthttp.StatusBadRequest, "", InvalidCredentials)
		return
	}

	randomToken := generateRandomToken(36)
	err = s.Tokens.Create(randomToken, user.ID)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusInternalServerError, err.Error(), "")
		return
	}

	response := authResponse{
		Auth:  *user,
		Token: randomToken,
	}

	okResponse(ctx, fasthttp.StatusOK, response)
}

func (s *Server) handleAuth(ctx *fasthttp.RequestCtx) {
	u, err := s.getUserByToken(ctx)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusUnauthorized, "Unauthorized", Unauthorized)
		return
	}

	okResponse(ctx, fasthttp.StatusOK, u)
}

func (s *Server) handleLogout(ctx *fasthttp.RequestCtx) {
	token := ctx.UserValue("token").(string)
	t, err := s.Tokens.GetByToken(token)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusUnauthorized, "Unauthorized", Unauthorized)
		return
	}

	s.Tokens.Delete(t.Token)
	okResponse(ctx, fasthttp.StatusOK, map[string]string{})
}

type settingsPassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (s *Server) ChangePassword(ctx *fasthttp.RequestCtx) {
	var sp settingsPassword
	json.Unmarshal(ctx.PostBody(), &sp)

	u, err := s.getUserByToken(ctx)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusUnauthorized, "Unauthorized", Unauthorized)
		return
	}

	if !checkPasswordHash(sp.OldPassword, u.Password) {
		errorResponse(ctx, fasthttp.StatusBadRequest, "", InvalidCredentials)
		return
	}

	if checkPasswordHash(sp.NewPassword, u.Password) {
		errorResponse(ctx, fasthttp.StatusBadRequest, "", PasswordEquals)
		return
	}

	passRegexp := regexp.MustCompile(`^.{6,}$`)
	if !passRegexp.MatchString(sp.NewPassword) {
		errorResponse(ctx, fasthttp.StatusBadRequest, "", InvalidPassword)
		return
	}

	hashPass, _ := hashPassword(sp.NewPassword)
	err = s.Users.EditPassword(u.ID, hashPass)
	if err != nil {
		errorResponse(ctx, fasthttp.StatusInternalServerError, err.Error(), "")
		return
	}

	okResponse(ctx, fasthttp.StatusCreated, map[string]string{})
}

func hashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return bytes, err
}

func checkPasswordHash(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}

// https://stackoverflow.com/a/31832326/21711976
func generateRandomToken(length int) string {
	var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_")
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
