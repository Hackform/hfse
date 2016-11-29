package pionen

import (
	"bytes"
	"crypto/rand"
	"errors"
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/route"
	"github.com/Hackform/hfse/route/liberty"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/scrypt"
	"net/http"
	"time"
)

type (
	Pionen struct {
		id         kappa.Const
		path       string
		routes     *route.RouteSubstrate
		userRoute  kappa.Const
		signingKey []byte
	}

	authClaim struct {
		jwt.StandardClaims
	}
)

const (
	salt_length     = 32
	hash_length     = 64
	work_factor     = 16384
	mem_blocksize   = 8
	parallel_factor = 1
)

func New(path string, routes *route.RouteSubstrate, userRoute kappa.Const) *Pionen {
	return &Pionen{
		path:      path,
		routes:    routes,
		userRoute: userRoute,
	}
}

func (p *Pionen) SetId(id kappa.Const) kappa.Const {
	p.id = id
	return p.id
}

func (p *Pionen) GetId() kappa.Const {
	return p.id
}

func (p *Pionen) GetPath() string {
	return p.path
}

//////////////
// Handlers //
//////////////

func (p *Pionen) VerifyUser(userid, password string) bool {
	userData := p.routes.Get(p.userRoute).(*liberty.Liberty)

	result := liberty.NewUser()
	done := userData.GetUser(userid, &result)
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
	if p.VerifyUser(userid, password) {
		claims := authClaim{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
				Issuer:    "hfse-pionen",
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString(p.signingKey)
	} else {
		return "", errors.New("not authorized")
	}
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
