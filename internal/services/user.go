package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"chatie.com/internal/domain"
	"chatie.com/utils"
	"github.com/sirupsen/logrus"
)

type userRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type UserService struct {
	repo    userRepository
	logger  logrus.FieldLogger
	timeout time.Duration
}

func NewUserService(logger logrus.FieldLogger, repo userRepository) *UserService {
	return &UserService{
		repo:    repo,
		logger:  logger,
		timeout: time.Duration(3) * time.Second,
	}
}

func (s *UserService) CreateUser(c context.Context, req *domain.CreateUserReq) (*domain.CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	res := &domain.CreateUserRes{
		ID:       strconv.Itoa(int(r.ID)),
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

func (s *UserService) Login(c context.Context, req *domain.LoginUserReq) (*domain.LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	fmt.Println("[user service login 1]", req.Email, req.Password, "[user service login 1]")
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &domain.LoginUserRes{}, err
	}
	fmt.Println("[user service login 2]", user)
	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		return &domain.LoginUserRes{}, err
	}

	token, err := utils.GenerateToken(int(user.ID), user.Username)
	if err != nil {
		return &domain.LoginUserRes{}, err
	}

	return &domain.LoginUserRes{ID: strconv.Itoa(int(user.ID)), Username: user.Username, AccessToken: token}, nil
}
