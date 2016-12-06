package pionenroute

import (
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/model/libertymodel"
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
		loginAttempt := new(RequestLogin)
		c.Bind(loginAttempt)
		if jwtString, err := auth.GetJWT(loginAttempt.Value.Id, loginAttempt.Value.Password); err == nil {
			return c.JSON(http.StatusOK, RequestJWT{Value: JWTToken{Token: jwtString}})
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid password")
		}
	})

	g.POST("/verify", func(c echo.Context) error {
		req := new(RequestJWT)
		c.Bind(req)
		if auth.VerifyJWT(req.Value.Token, access.USER) {
			return c.JSON(http.StatusOK, true)
		} else {
			return c.JSON(http.StatusBadRequest, false)
		}
	})

	// for testing
	g.POST("/unstable/decode", func(c echo.Context) error {
		req := new(RequestJWT)
		c.Bind(req)

		token, err := auth.ParseJWT(req.Value.Token)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid token")
		}

		return c.JSON(http.StatusOK, token)
	})
}

func (p *PionenRoute) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
