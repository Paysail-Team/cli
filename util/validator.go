/*
Copyright © 2023 Syro team <info@syro.com>
*/
package util

import (
	"net/mail"
)

func IsEmailAddressValid(emailAddress string) (isValid bool) {
	_, err := mail.ParseAddress(emailAddress)
	if err != nil {
		return false
	}
	return true
}
