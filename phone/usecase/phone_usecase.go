package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/haritsrizkall/jti-test/constant"
	"github.com/haritsrizkall/jti-test/domain"
	"github.com/haritsrizkall/jti-test/phone/websocket"
	"github.com/haritsrizkall/jti-test/utils"
)

type phoneUsecase struct {
	phoneHub        *websocket.Hub
	phoneRepository domain.PhoneRepository
}

func NewPhoneUsecase(phoneRepository domain.PhoneRepository, phoneHub *websocket.Hub) domain.PhoneUsecase {
	return &phoneUsecase{
		phoneRepository: phoneRepository,
		phoneHub:        phoneHub,
	}
}

func (u *phoneUsecase) GetAll(ctx context.Context) ([]domain.Phone, error) {
	phones, err := u.phoneRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// var oddPhones []domain.Phone
	// var evenPhones []domain.Phone

	// for _, phone := range phones {
	// 	phoneInt, _ := strconv.Atoi(phone.Number[len(phone.Number)-1:])
	// 	if phoneInt%2 == 0 {
	// 		evenPhones = append(evenPhones, phone)
	// 	} else {
	// 		oddPhones = append(oddPhones, phone)
	// 	}
	// }

	// response := &domain.GetAllPhoneResponse{
	// 	Odd:  oddPhones,
	// 	Even: evenPhones,
	// }

	return phones, nil
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
	if !utils.IsValidProvider(phone.Provider) && !utils.IsValidPhoneNumber(phone.Number) {
		return nil, errors.New(constant.ErrBadRequest)
	}

	existPhone, err := u.phoneRepository.GetByNumber(ctx, phone.Number)
	if err != nil && err.Error() != constant.ErrNoRowsInResultSet {
		return nil, err
	}
	if existPhone.Number != "" && existPhone.ID != phone.ID {
		return nil, errors.New(constant.ErrPhoneNumberExist)
	}

	err = u.phoneRepository.Update(ctx, phone)
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

	existPhone, err := u.phoneRepository.GetByNumber(ctx, phone.Number)
	if err != nil && err.Error() != constant.ErrNoRowsInResultSet {
		return nil, err
	}
	if existPhone.Number != "" {
		return nil, errors.New(constant.ErrPhoneNumberExist)
	}

	id, err := u.phoneRepository.Store(ctx, phone)
	if err != nil {
		return nil, err
	}

	phone.ID = id

	// broadcast to all client
	data := []domain.Phone{phone}
	dataJson, _ := json.Marshal(data)
	u.phoneHub.Broadcast <- dataJson

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

func (u *phoneUsecase) AutoGenerate(ctx context.Context) error {
	count := 25
	phones := utils.GeneratePhones(count)

	ids, err := u.phoneRepository.StoreBulk(ctx, phones)
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		phones[i].ID = ids[i]
	}

	// broadcast to all client
	data, _ := json.Marshal(phones)
	u.phoneHub.Broadcast <- data

	return nil
}
