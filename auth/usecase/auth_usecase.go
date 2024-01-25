package usecase

import (
	"context"
	"errors"
	"log"

	"github.com/haritsrizkall/jti-test/constant"
	"github.com/haritsrizkall/jti-test/domain"
	"github.com/haritsrizkall/jti-test/pkg"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepository domain.UserRepository
	googleOauth    pkg.GoogleOAuth
}

func NewAuthUsecase(userRepository domain.UserRepository, oauth pkg.GoogleOAuth) domain.AuthUsecase {
	return &authUsecase{
		userRepository: userRepository,
		googleOauth:    oauth,
	}
}

func (u *authUsecase) Register(ctx context.Context, request *domain.RegisterRequest) error {
	user, err := u.userRepository.GetByEmail(ctx, request.Email)
	if err != nil {
		if err.Error() != constant.ErrNoRowsInResultSet {
			return err
		}
	}

	if user.ID != 0 {
		return errors.New(constant.ErrEmailAlreadyExist)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user = domain.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
	}

	_, err = u.userRepository.Store(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *authUsecase) Login(ctx context.Context, request *domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := u.userRepository.GetByEmail(ctx, request.Email)
	if err != nil {
		if err.Error() == constant.ErrNoRowsInResultSet {
			return nil, errors.New(constant.ErrEmailNotFound)
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, errors.New(constant.ErrWrongPassword)
	}

	token, err := pkg.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	response := domain.LoginResponse{
		Token: token,
	}

	return &response, nil
}

func (u *authUsecase) LoginWithGoogle(ctx context.Context) string {
	url := u.googleOauth.GetAuthCodeURL("state")
	return url
}

func (u *authUsecase) LoginWithGoogleCallback(ctx context.Context, code string) (string, error) {
	userInfo, err := u.googleOauth.GetUserInfo(ctx, code)
	if err != nil {
		log.Println(err)
		return "", err
	}

	user, err := u.userRepository.GetByEmail(ctx, userInfo.Email)
	if err != nil {
		if err.Error() != constant.ErrNoRowsInResultSet {
			return "", err
		}

		user = domain.User{
			Name:  userInfo.Name,
			Email: userInfo.Email,
		}

		id, err := u.userRepository.Store(ctx, user)
		if err != nil {
			return "", err
		}

		token, err := pkg.GenerateToken(id)
		if err != nil {
			return "", err
		}

		return token, nil
	}

	token, err := pkg.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
