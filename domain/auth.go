package domain

import "context"

type AuthUsecase interface {
	Register(ctx context.Context, request *RegisterRequest) error
	Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error)
	LoginWithGoogle(ctx context.Context) string
	LoginWithGoogleCallback(ctx context.Context, code string) (string, error)
}

// Request
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
