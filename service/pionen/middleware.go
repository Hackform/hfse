package pionen

import (
	"github.com/Hackform/hfse/model/pionenmodel"
	"github.com/Hackform/hfse/service/pionen/access"
	"github.com/labstack/echo"
	"net/http"
)

type (
	UserIdFunc func(echo.Context) string
)

func (p *Pionen) MAuthLevel(level uint8) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := pionenmodel.GetRequestJWT(c)
			if p.VerifyJWTLevel(req.Value.Token, level) {
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
			req := pionenmodel.GetRequestJWT(c)
			if p.VerifyJWTUserId(req.Value.Token, f(c)) {
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
