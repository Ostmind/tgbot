package utils

import (
	"fmt"
	"regexp"
	"tgbot/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func ValidateLengthAndStrings(password string) (string, error) {
	re := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$`)

	if !re.MatchString(password) {
		return "", fmt.Errorf("error validating password %w", models.ErrPasswordNotSecured)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error securing password %w", err)
	}

	return string(hashedPassword), nil
}

func ComparePassword(userPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("error comparing password %w", err)
	}

	return nil
}
