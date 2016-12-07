package libertyroute

import (
	"fmt"
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/model/libertymodel"
	"github.com/Hackform/hfse/route"
	"github.com/Hackform/hfse/service/himeji"
	"github.com/Hackform/hfse/service/pionen"
	"github.com/Hackform/hfse/service/pionen/access"
	pionenhash "github.com/Hackform/hfse/service/pionen/hash"
	"github.com/labstack/echo"
	"net/http"
)

type (
	LibertyRoute struct {
		route.RouteBase
		repoService kappa.Const
		authService kappa.Const
	}
)

func New(path string, repoService kappa.Const, authService kappa.Const) *LibertyRoute {
	l := &LibertyRoute{
		repoService: repoService,
		authService: authService,
	}
	l.RouteBase.SetPath(path)
	return l
}

//////////////
// Register //
//////////////

func (l *LibertyRoute) Register(g *echo.Group) {
	repo := l.GetService(l.repoService).(*himeji.Himeji)
	auth := l.GetService(l.authService).(*pionen.Pionen)

	g.GET("/:userid", func(c echo.Context) error {
		result := new(himeji.Data)
		done := libertymodel.GetUser(repo, c.Param("userid"), result)
		if <-done {
			return c.JSON(http.StatusOK, libertymodel.RequestPublicUser{Value: result.Value.(libertymodel.ModelUser).PublicUser})
		} else {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user %s not found", c.Param("userid")))
		}
	})

	g.GET("/:userid/private", func(c echo.Context) error {
		result := new(himeji.Data)
		done := libertymodel.GetUser(repo, c.Param("userid"), result)
		if <-done {
			return c.JSON(http.StatusOK, libertymodel.RequestModelUser{Value: result.Value.(libertymodel.ModelUser)})
		} else {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user %s not found", c.Param("userid")))
		}
	}, auth.MAuthUserUrlParam("userid"))

	g.POST("", func(c echo.Context) error {
		user := libertymodel.GetRequestPostUser(c)
		if len(user.Value.Id) < 1 || len(user.Value.Name) < 1 || len(user.Value.Password) < 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "json malformed")
		}

		if <-libertymodel.GetUser(repo, user.Value.Id, new(himeji.Data)) {
			return echo.NewHTTPError(http.StatusBadRequest, "user exists")
		}

		hash, salt, err := pionenhash.Hash(user.Value.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid password")
		}

		usermodel := libertymodel.ModelUser{
			PublicUser: user.Value.PublicUser,
			PrivateUser: libertymodel.PrivateUser{
				UserPermissions: libertymodel.UserPermissions{
					AccessLevel: access.USER,
					AccessTags:  make([]uint8, 0),
				},
				Hash: hash,
				Salt: salt,
			},
		}
		done := libertymodel.StoreUser(repo, &himeji.Data{Value: usermodel})
		if <-done {
			return c.JSON(http.StatusCreated, libertymodel.RequestPublicUser{Value: user.Value.PublicUser})
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not store user")
		}
	})
}

func (l *LibertyRoute) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
