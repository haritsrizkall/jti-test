package utils

import (
	"math/rand"

	"github.com/haritsrizkall/jti-test/constant"
	"github.com/haritsrizkall/jti-test/domain"
)

func GeneratePhone() domain.Phone {
	phone := domain.Phone{}
	phone.Provider = constant.Providers[rand.Intn(len(constant.Providers))]
	phone.Number = "08"
	for i := 0; i < 10; i++ {
		phone.Number += string(rune(48 + rand.Intn(10)))
	}
	return phone
}

func GeneratePhones(count int) []domain.Phone {
	var phones []domain.Phone
	for i := 0; i < count; i++ {
		phones = append(phones, GeneratePhone())
	}
	return phones
}
