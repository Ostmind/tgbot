package helpers

import (
	"fmt"
	"github.com/Ostmind/tgbot/internal/todo/models"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(password string) (string, error) {
	hasLower := false
	hasUpper := false
	hasDigit := false

	for _, r := range password {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}

	if !(hasLower && hasUpper && hasDigit && len(password) > 8) {
		return "", models.ErrPasswordNotSecured
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error securing password %w", err)
	}

	return string(hashedPassword), nil
}

func ComparePassword(userPasswordCookie string, userPasswordDB string) error {
	err := bcrypt.CompareHashAndPassword([]byte(userPasswordCookie), []byte(userPasswordDB))
	if err != nil {
		return fmt.Errorf("error comparing password %w", err)
	}

	return nil
}
