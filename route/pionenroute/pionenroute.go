package pionenroute

import (
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/model/pionenmodel"
	"github.com/Hackform/hfse/route"
	"github.com/Hackform/hfse/service/pionen"
	"github.com/Hackform/hfse/service/pionen/access"
	"github.com/labstack/echo"
	"net/http"
)

type (
	PionenRoute struct {
		route.RouteBase
		authService kappa.Const
	}
)

func New(path string, authService kappa.Const) *PionenRoute {
	p := &PionenRoute{
		authService: authService,
	}
	p.RouteBase.SetPath(path)
	return p
}

//////////////
// Register //
//////////////

func (p *PionenRoute) Register(g *echo.Group) {
	auth := p.GetService(p.authService).(*pionen.Pionen)

	g.POST("/login", func(c echo.Context) error {
		loginAttempt := pionenmodel.GetRequestLogin(c)
		if jwtString, err := auth.GetJWT(loginAttempt.Value.Id, loginAttempt.Value.Password); err == nil {
			return c.JSON(http.StatusOK, pionenmodel.RequestJWT{Value: pionenmodel.JWTToken{Token: jwtString}})
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid password")
		}
	})

	g.POST("/verify", func(c echo.Context) error {
		req := pionenmodel.GetRequestJWT(c)
		if auth.VerifyJWTLevel(req.Value.Token, access.USER) {
			return c.JSON(http.StatusOK, true)
		} else {
			return c.JSON(http.StatusUnauthorized, "unauthorized")
		}
	})

	// for testing
	g.POST("/unstable/decode", func(c echo.Context) error {
		req := pionenmodel.GetRequestJWT(c)

		claims, err := auth.ParseJWT(req.Value.Token)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid token")
		}

		return c.JSON(http.StatusOK, claims)
	})
}

func (p *PionenRoute) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
