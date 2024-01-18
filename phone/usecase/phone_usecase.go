package usecase

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/haritsrizkall/jti-test/constant"
	"github.com/haritsrizkall/jti-test/domain"
	"github.com/haritsrizkall/jti-test/utils"
)

type phoneUsecase struct {
	phoneRepository domain.PhoneRepository
}

func NewPhoneUsecase(phoneRepository domain.PhoneRepository) domain.PhoneUsecase {
	return &phoneUsecase{
		phoneRepository: phoneRepository,
	}
}

func (u *phoneUsecase) GetAll(ctx context.Context) (*domain.GetAllPhoneResponse, error) {
	phones, err := u.phoneRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var oddPhones []domain.Phone
	var evenPhones []domain.Phone

	for _, phone := range phones {
		phoneInt, _ := strconv.Atoi(phone.Number[len(phone.Number)-1:])
		if phoneInt%2 == 0 {
			evenPhones = append(evenPhones, phone)
		} else {
			oddPhones = append(oddPhones, phone)
		}
	}

	response := &domain.GetAllPhoneResponse{
		Odd:  oddPhones,
		Even: evenPhones,
	}

	return response, nil
}

func (u *phoneUsecase) GetByID(ctx context.Context, id int) (*domain.Phone, error) {
	phone, err := u.phoneRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &phone, nil
}

func (u *phoneUsecase) Update(ctx context.Context, request domain.UpdatePhoneRequest) (*domain.Phone, error) {
	phone := domain.Phone{
		ID:       request.ID,
		Number:   request.Number,
		Provider: request.Provider,
	}

	// validate phone
	if !utils.IsValidProvider(phone.Provider) {
		return nil, errors.New(constant.ErrInvalidProvider)
	}
	regexp := regexp.MustCompile(`^08[0-9]{9,13}$`)
	if !regexp.MatchString(phone.Number) {
		fmt.Println(phone.Number)
		return nil, errors.New(constant.ErrInvalidPhoneNumber)
	}

	err := u.phoneRepository.Update(ctx, phone)
	if err != nil {
		return nil, err
	}

	return &phone, nil
}

func (u *phoneUsecase) Create(ctx context.Context, request domain.CreatePhoneRequest) (*domain.Phone, error) {
	// validate phone number
	phone := domain.Phone{
		Number:   request.Number,
		Provider: request.Provider,
	}

	// validate phone
	if !utils.IsValidProvider(phone.Provider) {
		return nil, errors.New(constant.ErrInvalidProvider)
	}
	regexp := regexp.MustCompile(`^08[0-9]{9,13}$`)
	if !regexp.MatchString(phone.Number) {
		fmt.Println(phone.Number)
		return nil, errors.New(constant.ErrInvalidPhoneNumber)
	}

	id, err := u.phoneRepository.Store(ctx, phone)
	if err != nil {
		return nil, err
	}

	phone.ID = id

	return &phone, nil
}

func (u *phoneUsecase) Delete(ctx context.Context, id int) (*domain.Phone, error) {
	err := u.phoneRepository.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	response := &domain.Phone{
		ID: id,
	}

	return response, nil
}
