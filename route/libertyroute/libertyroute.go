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
	"github.com/pborman/uuid"
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

	g.GET("/:username", func(c echo.Context) error {
		result := new(himeji.Data)
		done := libertymodel.GetUserByUsername(repo, c.Param("username"), result)
		if <-done {
			return c.JSON(http.StatusOK, libertymodel.RequestPublicUser{Value: result.Value.(libertymodel.ModelUser).PublicUser})
		} else {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user %s not found", c.Param("username")))
		}
	})

	g.GET("/:uid/private", func(c echo.Context) error {
		result := new(himeji.Data)
		done := libertymodel.GetUser(repo, c.Param("uid"), result)
		if <-done {
			return c.JSON(http.StatusOK, libertymodel.RequestUserInfo{Value: result.Value.(libertymodel.UserInfo)})
		} else {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user %s not found", c.Param("uid")))
		}
	}, auth.MAuthUserUrlParam("uid"))

	g.POST("", func(c echo.Context) error {
		user := libertymodel.GetRequestPostUser(c)
		if len(user.Value.Username) < 1 || len(user.Value.Password) < 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "json malformed")
		}

		if <-libertymodel.GetUserByUsername(repo, user.Value.Username, new(himeji.Data)) {
			return echo.NewHTTPError(http.StatusBadRequest, "user exists")
		}

		hash, salt, err := pionenhash.Hash(user.Value.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid password")
		}

		uid := uuid.New()

		usermodel := libertymodel.ModelUser{
			Uid: libertymodel.Uid{
				Id: uid,
			},
			UserInfo: user.Value.UserInfo,
			UserSecurity: libertymodel.UserSecurity{
				Hash: hash,
				Salt: salt,
			},
			UserPermissions: libertymodel.UserPermissions{
				AccessLevel: access.USER,
				AccessTags:  make([]uint8, 0),
			},
		}
		done := libertymodel.StoreUser(repo, &himeji.Data{Value: usermodel})
		if <-done {
			return c.JSON(http.StatusCreated, libertymodel.RequestPublicUser{Value: user.Value.PublicUser})
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not store user")
		}
	})

	adminGroup := g.Group("/admin", auth.MAuthAdmin())

	adminGroup.PUT("/:username", func(c echo.Context) error {
		return c.String(http.StatusOK, "put user")
	})
}

func (l *LibertyRoute) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
