package validation

import (
	"net/mail"
)

func IsEmail(str *string) bool {
	_, err := mail.ParseAddress(*str)
	return err == nil
}
