package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJWT(username string, accessTokenExpDate time.Duration, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(accessTokenExpDate).Unix(),
		Subject:   username,
	})

	return token.SignedString([]byte(jwtSecret))
}

func NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	rand.Read(b)

	hashedRefreshToken, _ := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)

	return hex.EncodeToString(hashedRefreshToken), nil
}

func Parse(accessToken string, jwtSecret string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing algorithm: %v", token.Header["alg"])
		}

		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	} else {
		return "", fmt.Errorf("Invalid Token ")
	}
}
