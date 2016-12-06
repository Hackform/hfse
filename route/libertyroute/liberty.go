package libertyroute

import (
	"fmt"
	"github.com/Hackform/hfse/kappa"
	model "github.com/Hackform/hfse/model/liberty"
	"github.com/Hackform/hfse/route"
	"github.com/Hackform/hfse/service/himeji"
	"github.com/Hackform/hfse/service/pionen/access"
	pionenhash "github.com/Hackform/hfse/service/pionen/hash"
	"github.com/labstack/echo"
	"net/http"
)

type (
	LibertyRoute struct {
		route.RouteBase
		repoService kappa.Const
	}

	RequestPublicUser struct {
		Value model.PublicUser `json:"data"`
	}

	RequestModelUser struct {
		Value model.ModelUser `json:"data"`
	}

	PostUser struct {
		model.PublicUser
		model.UserPassword
	}

	RequestPostUser struct {
		Value PostUser `json:"data"`
	}
)

func New(path string, repoService kappa.Const) *LibertyRoute {
	l := &LibertyRoute{
		repoService: repoService,
	}
	l.RouteBase.SetPath(path)
	return l
}

//////////////
// Register //
//////////////

func (l *LibertyRoute) Register(g *echo.Group) {
	repo := l.GetService(l.repoService).(*himeji.Himeji)

	g.GET("/:userid", func(c echo.Context) error {
		result := new(himeji.Data)
		done := model.GetUser(repo, c.Param("userid"), result)
		if <-done {
			return c.JSON(http.StatusOK, RequestPublicUser{Value: result.Value.(model.ModelUser).PublicUser})
		} else {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user %s not found", c.Param("userid")))
		}
	})

	g.POST("", func(c echo.Context) error {
		user := new(RequestPostUser)
		err := c.Bind(user)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "json malformed")
		}
		hash, salt, err := pionenhash.Hash(user.Value.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid password")
		}
		usermodel := model.ModelUser{
			PublicUser: user.Value.PublicUser,
			PrivateUser: model.PrivateUser{
				UserPermissions: model.UserPermissions{
					AccessLevel: access.USER,
					AccessTags:  make([]uint8, 0),
				},
				Hash: hash,
				Salt: salt,
			},
		}
		done := model.StoreUser(repo, &himeji.Data{Value: usermodel})
		if <-done {
			return c.JSON(http.StatusCreated, RequestPublicUser{user.Value.PublicUser})
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not store user")
		}
	})
}

func (l *LibertyRoute) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
