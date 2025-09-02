package helpers

import (
	"net/http"
	"time"
)

func SetCookie(name string, value string, expDate time.Duration, httpOnly bool) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Expires = time.Now().Add(expDate)
	cookie.HttpOnly = httpOnly

	return cookie
}
