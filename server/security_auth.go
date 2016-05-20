package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

func getSessionIdFromBearerToken(token string) (int, error) {
	prefix := "sid-"
	if !strings.HasPrefix(token, prefix) {
		return 0, echo.NewHTTPError(http.StatusUnauthorized)
	}

	sidString := token[len(prefix):]
	i, err := strconv.ParseInt(sidString, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

func GetClientPubkey() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header().Get("Authorization")

			prefixLength := len("Bearer")

			if len(auth) <= len("Bearer")+1 || auth[:prefixLength] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			token := auth[prefixLength+1:]
			sessionId, err := getSessionIdFromBearerToken(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			c.Set("session-id", sessionId)
			return next(c)
		}
	}
}
