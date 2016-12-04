package pionen

import (
	"bytes"
	"crypto/rand"
	"errors"
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/route"
	"github.com/Hackform/hfse/route/liberty"
	"github.com/Hackform/hfse/service/himeji"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/scrypt"
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
)

const (
	salt_length     = 32
	hash_length     = 64
	work_factor     = 16384
	mem_blocksize   = 8
	parallel_factor = 1

	jwt_iss   = "hfse-pionen"
	jwt_hours = 72
)

func New(path string, routes *route.RouteSubstrate, userRoute kappa.Const) *Pionen {
	p := &Pionen{
		userRoute: userRoute,
	}
	p.RouteBase.SetPath(path)
	return p
}

//////////////
// Handlers //
//////////////

func (p *Pionen) VerifyUser(userid, password string) bool {
	userData := p.GetRoute(p.userRoute).(*liberty.Liberty)

	result := new(himeji.Data)
	done := userData.GetUser(userid, result)
	if !<-done {
		return false
	}
	user := result.Value.(liberty.ModelUser)
	dk, err := scrypt.Key([]byte(password), user.Salt, work_factor, mem_blocksize, parallel_factor, hash_length)
	if err != nil {
		return false
	}
	return bytes.Equal(dk, user.Hash)
}

func (p *Pionen) GetJWT(userid, password string) (string, error) {
	userData := p.GetRoute(p.userRoute).(*liberty.Liberty)

	result := new(himeji.Data)
	done := userData.GetUser(userid, result)
	if !<-done {
		return "", errors.New("cannot fetch user")
	}
	user := result.Value.(liberty.ModelUser)

	if p.VerifyUser(userid, password) {
		claims := authClaim{
			Issuer:    jwt_iss,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * jwt_hours).Unix(),
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
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	token, err := jwt.ParseWithClaims(tokenString, &authClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return p.signingKey, nil
	})

	if claims, ok := token.Claims.(*authClaim); ok && token.Valid && err != nil {
		userData := p.GetRoute(p.userRoute).(*liberty.Liberty)

		result := new(himeji.Data)
		done := userData.GetUser(claims.Id, result)
		if !<-done {
			return "", errors.New("cannot fetch user")
		}
		user := result.Value.(liberty.ModelUser)

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

func Hash(password string) (h, s []byte, e error) {
	salt := make([]byte, salt_length)
	_, err := rand.Read(salt)
	if err != nil {
		return make([]byte, hash_length), salt, err
	}
	hash, err := scrypt.Key([]byte(password), salt, work_factor, mem_blocksize, parallel_factor, hash_length)
	return hash, salt, err
}

//////////////
// Register //
//////////////

func (p *Pionen) Register(g *echo.Group) {
	g.POST("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "retrieve")
	})
}

func (p *Pionen) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
