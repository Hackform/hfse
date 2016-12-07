package pionenmodel

import (
	"github.com/Hackform/hfse/model/libertymodel"
	"github.com/labstack/echo"
)

type (
	Login struct {
		libertymodel.UserId
		libertymodel.UserPassword
	}

	RequestLogin struct {
		Value Login `json:"auth"`
	}

	JWTToken struct {
		Token string `json:"token"`
	}

	RequestJWT struct {
		Value JWTToken `json:"auth"`
	}
)

func GetRequestJWT(c echo.Context) *RequestJWT {
	req := new(RequestJWT)
	c.Bind(req)
	return req
}

func GetRequestLogin(c echo.Context) *RequestLogin {
	req := new(RequestLogin)
	c.Bind(req)
	return req
}
