package liberty

import (
	"fmt"
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/route"
	"github.com/Hackform/hfse/route/pionen"
	"github.com/Hackform/hfse/route/pionen/access"
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

	PublicUser struct {
		Id   string `json:"id" bson:"_id"`
		Name string `json:"name" bson:"name"`
	}

	PrivateUser struct {
		AccessLevel uint8   `json:"accesslevel" bson:"accesslevel"`
		AccessTags  []uint8 `json:"accesstags" bson:"accesstags"`
		Hash        []byte  `json:"hash" bson:"hash"`
		Salt        []byte  `json:"salt" bson:"salt"`
	}

	ModelUser struct {
		PublicUser
		PrivateUser
	}

	RequestPublicUser struct {
		Value PublicUser `json:"data" bson:"data"`
	}

	RequestModelUser struct {
		Value ModelUser `json:"data" bson:"data"`
	}

	PostUser struct {
		PublicUser
		Password string `json:"password" bson:"password"`
	}

	RequestPostUser struct {
		Value PostUser `json:"data" bson:"data"`
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
		hash, salt, err := pionen.Hash(user.Value.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid password")
		}
		usermodel := ModelUser{user.Value.PublicUser, PrivateUser{
			AccessLevel: access.USER,
			AccessTags:  make([]uint8, 0),
			Hash:        hash,
			Salt:        salt,
		}}
		done := l.StoreUser(&himeji.Data{Value: usermodel})
		if <-done {
			return c.JSON(http.StatusCreated, user)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not store user")
		}
	})
}

func (l *Liberty) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
