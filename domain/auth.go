package domain

import "context"

type AuthUsecase interface {
	LoginWithGoogle(ctx context.Context) string
	LoginWithGoogleCallback(ctx context.Context, code string) (string, error)
}
