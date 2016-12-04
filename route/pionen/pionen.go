package pionen

import (
	"errors"
	"fmt"
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/route"
	"github.com/Hackform/hfse/route/liberty"
	"github.com/Hackform/hfse/route/pionen/hash"
	"github.com/Hackform/hfse/service/himeji"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

type (
	Pionen struct {
		route.RouteBase
		userRoute  kappa.Const
		signingKey []byte
	}

	authClaim struct {
		jwt.StandardClaims
		liberty.PublicUser
		liberty.UserPermissions
	}

	Login struct {
		liberty.UserId
		liberty.UserPassword
	}

	RequestLogin struct {
		Data Login `json:"data"`
	}

	JWTToken struct {
		Token string `json:"token"`
	}

	RequestJWT struct {
		Data JWTToken `json:"data"`
	}
)

const (
	jwt_iss   = "hfse-pionen"
	jwt_hours = 72
)

func New(path string, userRoute kappa.Const) *Pionen {
	p := &Pionen{
		userRoute: userRoute,
	}
	p.RouteBase.SetPath(path)
	return p
}

//////////////
// Handlers //
//////////////

func (p *Pionen) VerifyUser(userid, password string) (bool, *liberty.ModelUser) {
	userData := p.GetRoute(p.userRoute).(*liberty.Liberty)

	result := new(himeji.Data)
	done := userData.GetUser(userid, result)
	if !<-done {
		return false, nil
	}
	user := result.Value.(liberty.ModelUser)

	if hash.Verify(password, user.Salt, user.Hash) {
		return true, &user
	}

	return false, nil
}

func (p *Pionen) GetJWT(userid, password string) (string, error) {
	if verify, user := p.VerifyUser(userid, password); verify {
		claims := authClaim{
			jwt.StandardClaims{
				Issuer:    jwt_iss,
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Hour * jwt_hours).Unix(),
			},
			user.PublicUser,
			user.PrivateUser.UserPermissions,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString(p.signingKey)
	} else {
		return "", errors.New("not authorized")
	}
}

func (p *Pionen) VerifyJWT(tokenString string, access uint8) bool {
	token, err := jwt.ParseWithClaims(tokenString, &authClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return p.signingKey, nil
	})

	if claims, ok := token.Claims.(*authClaim); ok && token.Valid && err != nil {
		if claims.AccessLevel <= access {
			return true
		}

		for _, i := range claims.AccessTags {
			if i == access {
				return true
			}
		}
	}

	return false
}

//////////////
// Register //
//////////////

func (p *Pionen) Register(g *echo.Group) {
	g.POST("/login", func(c echo.Context) error {
		loginAttempt := new(RequestLogin)
		c.Bind(loginAttempt)
		if jwtString, err := p.GetJWT(loginAttempt.Data.Id, loginAttempt.Data.Password); err != nil {
			return c.JSON(http.StatusOK, RequestJWT{Data: JWTToken{Token: jwtString}})
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid password")
		}
	})

	g.POST("/verify", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "verify")
	})
}

func (p *Pionen) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
