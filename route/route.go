package route

import (
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/service"
	"github.com/labstack/echo"
)

type (
	Route interface {
		SetId(k kappa.Const) kappa.Const
		GetId() kappa.Const
		GetPath() string
		SetPath(string) string
		SetServiceSubstrate(*service.ServiceSubstrate)
		GetService(kappa.Const) service.Service
		Register(*echo.Group)
		Middleware() []echo.MiddlewareFunc
	}

	RouteBase struct {
		service.ServiceBase
		path string
	}
)

func (r *RouteBase) GetPath() string {
	return r.path
}

func (r *RouteBase) SetPath(path string) string {
	r.path = path
	return path
}
