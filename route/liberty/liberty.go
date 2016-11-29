package liberty

import (
	"fmt"
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/service"
	"github.com/Hackform/hfse/service/himeji"
	"github.com/labstack/echo"
	"net/http"
)

type (
	Liberty struct {
		path        string
		substrate   *service.ServiceSubstrate
		repoService kappa.Const
	}

	modelUser struct {
		Id   string `json:"id" bson:"_id"`
		Name string `json:"name" bson:"name"`
	}
)

func New(path string, substrate *service.ServiceSubstrate, repoService kappa.Const) *Liberty {
	return &Liberty{
		path:        path,
		substrate:   substrate,
		repoService: repoService,
	}
}

func (l *Liberty) GetPath() string {
	return l.path
}

func (l *Liberty) Register(g *echo.Group) {
	collection := "Users"
	repo := l.substrate.Get(l.repoService).(*himeji.Himeji)

	g.GET("/:userid", func(c echo.Context) error {
		result := himeji.Data{
			Value: modelUser{},
		}
		done := repo.QueryId(collection, c.Param("userid"), &result)
		if <-done {
			return c.JSON(http.StatusOK, &result)
		} else {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user %s not found", c.Param("userid")))
		}
	})

	g.POST("", func(c echo.Context) error {
		user := himeji.Data{
			Value: modelUser{},
		}
		err := c.Bind(&user)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "json malformed")
		}
		done := repo.Insert(collection, &user)
		if <-done {
			return c.JSON(http.StatusCreated, &user)
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, "json malformed")
		}
	})
}

func (l *Liberty) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
