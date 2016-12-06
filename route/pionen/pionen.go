package pionen

import (
	"errors"
	"fmt"
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/route"
	"github.com/Hackform/hfse/route/liberty"
	"github.com/Hackform/hfse/route/pionen/access"
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
		jwt_iss    string
		jwt_hours  int
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
		Value Login `json:"data"`
	}

	JWTToken struct {
		Token string `json:"token"`
	}

	RequestJWT struct {
		Value JWTToken `json:"data"`
	}
)

func New(path string, signingKey, issuer string, hours int, userRoute kappa.Const) *Pionen {
	p := &Pionen{
		userRoute:  userRoute,
		signingKey: []byte(signingKey),
		jwt_iss:    issuer,
		jwt_hours:  hours,
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
				Issuer:    p.jwt_iss,
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Hour * time.Duration(p.jwt_hours)).Unix(),
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

func (p *Pionen) ParseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &authClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return p.signingKey, nil
	})
}

func (p *Pionen) VerifyJWT(tokenString string, level uint8) bool {
	token, err := p.ParseJWT(tokenString)

	if claims, ok := token.Claims.(*authClaim); err == nil && ok && token.Valid && claims.VerifyIssuer(p.jwt_iss, true) {
		if level > access.DIVIDE {
			for _, i := range claims.AccessTags {
				if i == level {
					return true
				}
			}
		} else {
			if claims.AccessLevel <= level {
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
		if jwtString, err := p.GetJWT(loginAttempt.Value.Id, loginAttempt.Value.Password); err == nil {
			return c.JSON(http.StatusOK, RequestJWT{Value: JWTToken{Token: jwtString}})
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid password")
		}
	})

	// for testing
	g.POST("/decode", func(c echo.Context) error {
		req := new(RequestJWT)
		c.Bind(req)

		token, err := p.ParseJWT(req.Value.Token)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid token")
		}

		return c.JSON(http.StatusOK, token)
	})

	g.POST("/verify", func(c echo.Context) error {
		req := new(RequestJWT)
		c.Bind(req)
		if p.VerifyJWT(req.Value.Token, access.USER) {
			return c.JSON(http.StatusOK, true)
		} else {
			return c.JSON(http.StatusBadRequest, false)
		}
	})
}

func (p *Pionen) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
