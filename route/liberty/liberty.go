package liberty

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type (
	Liberty struct {
		path string
	}
)

func New(path string) *Liberty {
	return &Liberty{
		path: path,
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
