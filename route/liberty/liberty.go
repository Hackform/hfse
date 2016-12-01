package liberty

import (
	"fmt"
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/route/pionen/access"
	"github.com/Hackform/hfse/service"
	"github.com/Hackform/hfse/service/himeji"
	"github.com/labstack/echo"
	"net/http"
)

type (
	Liberty struct {
		id          kappa.Const
		path        string
		services    *service.ServiceSubstrate
		repoService kappa.Const
	}

	ModelUser struct {
		Id          string  `json:"id" bson:"_id"`
		Name        string  `json:"name" bson:"name"`
		AccessLevel uint8   `json:"accesslevel" bson:"accesslevel"`
		AccessTags  []uint8 `json:"accesstags" bson:"accesstags"`
		Hash        []byte  `json:"hash" bson:"hash"`
		Salt        []byte  `json:"salt" bson:"salt"`
	}

	JsonUser struct {
		Value ModelUser `json:"data" bson:"data"`
	}
)

var (
	collection = "Users"
)

func New(path string, services *service.ServiceSubstrate, repoService kappa.Const) *Liberty {
	return &Liberty{
		path:        path,
		services:    services,
		repoService: repoService,
	}
}

func (l *Liberty) SetId(k kappa.Const) kappa.Const {
	l.id = k
	return k
}

func (l *Liberty) GetId() kappa.Const {
	return l.id
}

func (l *Liberty) GetPath() string {
	return l.path
}

//////////////
// Handlers //
//////////////

func (l *Liberty) GetUser(userid string, result *himeji.Data) <-chan bool {
	repo := l.services.Get(l.repoService).(*himeji.Himeji)
	return repo.QueryId(collection, userid, result)
}

func (l *Liberty) StoreUser(user *himeji.Data) <-chan bool {
	repo := l.services.Get(l.repoService).(*himeji.Himeji)
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
			return c.JSON(http.StatusOK, result)
		} else {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user %s not found", c.Param("userid")))
		}
	})

	g.POST("", func(c echo.Context) error {
		user := new(JsonUser)
		err := c.Bind(user)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "json malformed")
		}
		user.Value.AccessLevel = access.USER
		user.Value.AccessTags = make([]uint8, 0)
		done := l.StoreUser(&himeji.Data{Value: user.Value})
		if <-done {
			return c.JSON(http.StatusCreated, user)
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, "json malformed")
		}
	})
}

func (l *Liberty) Middleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}
