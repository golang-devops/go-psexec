package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func GetClientPubkey() echo.HandlerFunc {
	return func(c *echo.Context) error {
		auth := c.Request().Header.Get("Authorization")

		prefixLength := len("Bearer")

		if len(auth) <= len("Bearer")+1 || auth[:prefixLength] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		token := auth[prefixLength+1:]
		clientPubKey, err := getClientPubkeyFromToken(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		c.Set("session-token", clientPubKey)
		return nil
	}
}
