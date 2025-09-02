package auth

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type UserSession struct {
	Authenticated bool
	TwoFACode     string
	CodeExpire    time.Time
}

type AuthManager struct {
	sessions map[int64]*UserSession
	mu       sync.Mutex
}

func NewAuthManager() *AuthManager {
	return &AuthManager{
		sessions: make(map[int64]*UserSession),
	}
}

func (a *AuthManager) GenerateCode(chatID int64) string {
	a.mu.Lock()
	defer a.mu.Unlock()

	code := generateCode()
	a.sessions[chatID] = &UserSession{
		Authenticated: false,
		TwoFACode:     code,
		CodeExpire:    time.Now().Add(5 * time.Minute),
	}
	return code
}

func (a *AuthManager) VerifyCode(chatID int64, code string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	session, exists := a.sessions[chatID]
	if !exists {
		return false
	}
	if session.CodeExpire.Before(time.Now()) {
		delete(a.sessions, chatID)
		return false
	}
	if session.TwoFACode == code {
		session.Authenticated = true
		session.TwoFACode = ""
		return true
	}
	return false
}

func (a *AuthManager) IsAuthenticated(chatID int64) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	session, exists := a.sessions[chatID]
	return exists && session.Authenticated
}

func generateCode() string {
	return strconv.Itoa(rand.Intn(900000) + 100000) // 6 digit
}
