package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"chatie.com/internal/domain"
	"chatie.com/utils"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

const (
	WrongCredentials = "wrong credentials"
	WrongPassword    = "wrong password"
	WrongEmail       = "wrong email"
	UnexpectedError  = "unexpected error"
)

type UserService interface {
	CreateUser(c context.Context, req *domain.CreateUserReq) (*domain.CreateUserRes, error)
	Login(c context.Context, req *domain.LoginUserReq) (*domain.LoginUserRes, error)
}

type UserHandler struct {
	logger  logrus.FieldLogger
	service UserService
}

func NewUserHandler(logger logrus.FieldLogger, service UserService) *UserHandler {
	return &UserHandler{
		logger:  logger,
		service: service,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, UnexpectedError)
		return
	}

	passwordOk := utils.IsPasswordComplex(user.Password)
	if !passwordOk {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, WrongPassword)
		return
	}

	emailOk := utils.IsMail(user.Email)
	if !emailOk {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, WrongEmail)
		return
	}

	res, err := h.service.CreateUser(r.Context(), &user)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user domain.LoginUserReq
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, err.Error())
		return
	}

	logedUser, err := h.service.Login(r.Context(), &user)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, WrongCredentials)
		return
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    logedUser.AccessToken,
		Expires:  time.Now().Add(time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, logedUser)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-1000 * time.Hour),
		Path:     "/",
		HttpOnly: false,
	}

	http.SetCookie(w, &cookie)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, "success!")
}
