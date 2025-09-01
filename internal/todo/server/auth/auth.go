package auth

import (
	"fmt"
	"github.com/Ostmind/tgbot/internal/todo/config"
	helpers2 "github.com/Ostmind/tgbot/internal/todo/helpers"
	"github.com/Ostmind/tgbot/internal/todo/server/handler/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Authentication(echo echo.Context, manager *user.Controller, cfgAuth config.AuthConfig) (err error) {
	jwtTokenCookie, _ := echo.Cookie("AccessToken")

	refreshTokenCookie, _ := echo.Cookie("RefreshToken")

	userNameCookie, err := echo.Cookie("username")
	if err != nil {
		return echo.NoContent(http.StatusBadRequest)
	}

	passwordCookie, err := echo.Cookie("password")
	if err != nil {
		return echo.NoContent(http.StatusBadRequest)
	}

	userRepository, err := manager.Manager.GetUserByTelegramID(echo.Request().Context(), userNameCookie.Value)
	if err != nil {
		return err
	}

	if jwtTokenCookie != nil {
		jwtUsername, err := helpers2.Parse(jwtTokenCookie.Value, cfgAuth.JWTSecret)
		if err != nil || jwtUsername != userNameCookie.Value {
			return echo.NoContent(http.StatusBadRequest)
		}

		return nil
	}

	if refreshTokenCookie != nil {
		err = helpers2.ComparePassword(userRepository.RefreshToken, refreshTokenCookie.Value)
		if err != nil {
			return echo.NoContent(http.StatusBadRequest)
		}
	}

	err = helpers2.ComparePassword(userRepository.Password, passwordCookie.Value)
	if err != nil {
		return fmt.Errorf("error while authentication %s", err)
	}

	jwt, err := helpers2.GenerateJWT(userRepository.UserName, cfgAuth.JWTAccessTokenTTL, cfgAuth.JWTSecret)
	if err != nil {
		return echo.NoContent(http.StatusInternalServerError)
	}

	refreshToken, err := helpers2.NewRefreshToken()
	if err != nil {
		return echo.NoContent(http.StatusInternalServerError)
	}

	cookie := helpers2.SetCookie("AccessToken", jwt, cfgAuth.JWTAccessTokenTTL, false)

	echo.SetCookie(cookie)

	cookie = helpers2.SetCookie("RefreshToken", refreshToken, cfgAuth.JWTRefreshTokenTTL, true)

	echo.SetCookie(cookie)

	return nil
}
