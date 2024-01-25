package domain

import "context"

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (User, error)
	Store(ctx context.Context, user User) (int, error)
}
