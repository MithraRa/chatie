package queries

import (
	"context"

	"chatie.com/internal/domain"
	"chatie.com/internal/repository"
)

const insertUserQuery = `INSERT INTO users(username, password, email) VALUES ($1, $2, $3) returning id`

func (q *Queries) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	var lastInsertId int
	err := q.pool.QueryRow(ctx, insertUserQuery, user.Username, user.Password, user.Email).Scan(&lastInsertId)
	if err != nil {
		return &domain.User{}, repository.ErrPoolEmpty
	}

	user.ID = uint64(lastInsertId)

	return user, nil
}

const selectUserByEmail = `SELECT id, email, username, password FROM users WHERE email = $1`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := q.pool.QueryRow(ctx, selectUserByEmail, email).Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return &domain.User{}, nil
	}

	return &user, nil
}
