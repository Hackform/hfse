package liberty

import (
	"fmt"
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/route"
	"github.com/Hackform/hfse/route/pionen/access"
	pionenhash "github.com/Hackform/hfse/route/pionen/hash"
	"github.com/Hackform/hfse/service/himeji"
	"github.com/labstack/echo"
	"net/http"
)

type (
	Liberty struct {
		route.RouteBase
		repoService kappa.Const
	}

	////////////////
	// User Model //
	////////////////

	UserId struct {
		Id string `json:"id" bson:"_id"`
	}

	PublicUser struct {
		UserId
		Name string `json:"name" bson:"name"`
	}

	UserPermissions struct {
		AccessLevel uint8   `json:"accesslevel" bson:"accesslevel"`
		AccessTags  []uint8 `json:"accesstags" bson:"accesstags"`
	}

	PrivateUser struct {
		UserPermissions
		Hash []byte `json:"hash" bson:"hash"`
		Salt []byte `json:"salt" bson:"salt"`
	}

	ModelUser struct {
		PublicUser
		PrivateUser
	}

	RequestPublicUser struct {
		Value PublicUser `json:"data"`
	}

	RequestModelUser struct {
		Value ModelUser `json:"data"`
	}

	UserPassword struct {
		Password string `json:"password"`
	}

	PostUser struct {
		PublicUser
		UserPassword
	}

	RequestPostUser struct {
		Value PostUser `json:"data"`
	}
)

var (
	collection = "Users"
)

func New(path string, repoService kappa.Const) *Liberty {
	l := &Liberty{
		repoService: repoService,
	}
	l.RouteBase.SetPath(path)
	return l
}

//////////////
// Handlers //
//////////////

func (l *Liberty) GetUser(userid string, result *himeji.Data) <-chan bool {
	repo := l.GetService(l.repoService).(*himeji.Himeji)
	return repo.QueryId(collection, userid, result)
}

func (l *Liberty) StoreUser(user *himeji.Data) <-chan bool {
	repo := l.GetService(l.repoService).(*himeji.Himeji)
	return repo.Insert(collection, user)
}

//////////////
// Register //
//////////////

func (l *Liberty) Register(g *echo.Group) {
	g.GET("/:userid", func(c echo.Context) error {
		result := new(himeji.Data)
		done := l.GetUser(c.Param("userid"), result)
		if <-done {
			return c.JSON(http.StatusOK, RequestPublicUser{Value: result.Value.(ModelUser).PublicUser})
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
		usermodel := ModelUser{user.Value.PublicUser, PrivateUser{
			UserPermissions{
				AccessLevel: access.USER,
				AccessTags:  make([]uint8, 0),
			},
			hash,
			salt,
		}}
		done := l.StoreUser(&himeji.Data{Value: usermodel})
		if <-done {
			return c.JSON(http.StatusCreated, user.Value.PublicUser)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not store user")
		}
	})
}

func (l *Liberty) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
