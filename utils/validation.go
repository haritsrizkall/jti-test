package utils

import (
	"regexp"

	"github.com/haritsrizkall/jti-test/constant"
)

func IsValidProvider(provider string) bool {
	for _, validProvider := range constant.Providers {
		if validProvider == provider {
			return true
		}
	}
	return false
}

func IsValidPhoneNumber(number string) bool {
	regexp := regexp.MustCompile(`^08[0-9]{9,13}$`)
	if !regexp.MatchString(number) {
		return false
	}
	return true
}
