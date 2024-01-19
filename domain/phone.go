package domain

import "context"

type Phone struct {
	ID       int    `json:"id"`
	Number   string `json:"number"`
	Provider string `json:"provider"`
}

type PhoneRepository interface {
	GetAll(ctx context.Context) ([]Phone, error)
	GetByID(ctx context.Context, id int) (Phone, error)
	Update(ctx context.Context, phone Phone) error
	Store(ctx context.Context, phone Phone) (int, error)
	Delete(ctx context.Context, id int) error
}

type PhoneUsecase interface {
	GetAll(ctx context.Context) (*GetAllPhoneResponse, error)
	GetByID(ctx context.Context, id int) (*Phone, error)
	Update(ctx context.Context, request UpdatePhoneRequest) (*Phone, error)
	Create(ctx context.Context, request CreatePhoneRequest) (*Phone, error)
	Delete(ctx context.Context, id int) (*Phone, error)
}

// request
type CreatePhoneRequest struct {
	Number   string `json:"number" validate:"required"`
	Provider string `json:"provider" validate:"required"`
}

type UpdatePhoneRequest struct {
	ID       int    `json:"id" validate:"required"`
	Number   string `json:"number" validate:"required"`
	Provider string `json:"provider" validate:"required"`
}

// response
type GetAllPhoneResponse struct {
	Odd  []Phone `json:"odd"`
	Even []Phone `json:"even"`
}
