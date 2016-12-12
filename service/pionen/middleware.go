package pionen

import (
	"errors"
	"github.com/Hackform/hfse/service/pionen/access"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

type (
	UserIdFunc func(echo.Context) string
)

func GetAuthHeaderToken(c echo.Context) (string, error) {
	req := strings.Fields(c.Request().Header.Get("Authorization"))
	if len(req) != 2 || req[0] != "Bearer" {
		return "", errors.New("Not authorized")
	}
	return req[1], nil
}

func (p *Pionen) MAuthLevel(level uint8) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := GetAuthHeaderToken(c)
			if err == nil && p.VerifyJWTLevel(token, level) {
				return next(c)
			} else {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}
		}
	}
}

func (p *Pionen) MAuthRoot() echo.MiddlewareFunc {
	return p.MAuthLevel(access.ROOT)
}

func (p *Pionen) MAuthAdmin() echo.MiddlewareFunc {
	return p.MAuthLevel(access.ADMIN)
}

func (p *Pionen) MAuthMod() echo.MiddlewareFunc {
	return p.MAuthLevel(access.MOD)
}

func (p *Pionen) MAuthUser() echo.MiddlewareFunc {
	return p.MAuthLevel(access.USER)
}

func (p *Pionen) MAuthUserId(f UserIdFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := GetAuthHeaderToken(c)
			if err == nil && p.VerifyJWTUsername(token, f(c)) {
				return next(c)
			} else {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}
		}
	}
}

func (p *Pionen) MAuthUserUrlParam(param string) echo.MiddlewareFunc {
	return p.MAuthUserId(func(c echo.Context) string {
		return c.Param(param)
	})
}
