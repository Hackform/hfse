package pionen

import (
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/route"
	"github.com/Hackform/hfse/route/liberty"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
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

func (p *Pionen) GetJWT() (string, error) {
	userData := p.routes.Get(p.userRoute).(*liberty.Liberty)
	claims := authClaim{
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "pionen",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.signingKey)
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
