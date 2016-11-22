package liberty

import (
	"fmt"
	"net/http"

	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/service"
	"github.com/labstack/echo"
)

type (
	Liberty struct {
		path        string
		substrate   *service.ServiceSubstrate
		repoService kappa.Const
	}
)

func New(path string, substrate *service.ServiceSubstrate, repoService kappa.Const) *Liberty {
	return &Liberty{
		path:        path,
		substrate:   substrate,
		repoService: repoService,
	}
}

func (l *Liberty) Register(g *echo.Group) {
	g.GET("/:userid", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("fetching %s\n", c.Param("userid")))
	})
}

func (l *Liberty) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
