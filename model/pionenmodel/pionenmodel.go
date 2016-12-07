package pionenmodel

import (
	"github.com/Hackform/hfse/model/libertymodel"
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
