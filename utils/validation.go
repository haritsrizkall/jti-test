package utils

import "github.com/haritsrizkall/jti-test/constant"

func IsValidProvider(provider string) bool {
	for _, validProvider := range constant.Providers {
		if validProvider == provider {
			return true
		}
	}
	return false
}
