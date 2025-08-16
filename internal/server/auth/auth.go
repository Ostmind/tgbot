package auth

import (
	"net/http"
	"tgbot/internal/config"
	"tgbot/internal/server/handler/user"
	"tgbot/internal/utils"

	"github.com/labstack/echo/v4"
)

func Authentication(echo echo.Context, manager *user.Controller, cfgAuth config.AuthConfig) error {
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
		jwtUsername, err := utils.Parse(jwtTokenCookie.Value, cfgAuth.JWTSecret)
		if err != nil || jwtUsername != userNameCookie.Value {
			return echo.NoContent(http.StatusBadRequest)
		}

		return nil // go out of auth then JWT is Valid
	}

	//check if we have refresh token and it's valid
	if refreshTokenCookie != nil {
		err = utils.ComparePassword(userRepository.RefreshToken, refreshTokenCookie.Value)
		if err != nil {
			return echo.NoContent(http.StatusBadRequest)
		}
	}

	err = utils.ComparePassword(userRepository.Password, passwordCookie.Value)
	if err != nil {
		return err
	}

	jwt, err := utils.GenerateJWT(userRepository.UserName, cfgAuth.JWTAccessTokenTTL, cfgAuth.JWTSecret)
	if err != nil {
		return echo.NoContent(http.StatusInternalServerError)
	}

	refreshToken, err := utils.NewRefreshToken()
	if err != nil {
		return echo.NoContent(http.StatusInternalServerError)
	}

	cookie := utils.SetCookie("AccessToken", jwt, cfgAuth.JWTAccessTokenTTL, false)

	echo.SetCookie(cookie)

	cookie = utils.SetCookie("RefreshToken", refreshToken, cfgAuth.JWTRefreshTokenTTL, true)

	echo.SetCookie(cookie)

	return nil
}
