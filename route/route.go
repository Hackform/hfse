package route

import (
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/service"
	"github.com/labstack/echo"
)

type (
	RouteSubstrate struct {
		routeKappa *kappa.Kappa
		routes     map[kappa.Const]Route
	}

	Route interface {
		SetId(k kappa.Const) kappa.Const
		GetId() kappa.Const
		GetPath() string
		SetPath(string) string
		Register(*echo.Group)
		Middleware() []echo.MiddlewareFunc
	}

	RouteBase struct {
		service.ServiceBase
		path      string
		substrate *RouteSubstrate
	}
)

func New() *RouteSubstrate {
	return &RouteSubstrate{
		routeKappa: kappa.New(),
		routes:     make(map[kappa.Const]Route),
	}
}

func (r *RouteSubstrate) Set(rou Route) kappa.Const {
	k := r.routeKappa.Get()
	rou.SetId(k)
	r.routes[k] = rou
	return k
}

func (r *RouteSubstrate) Get(k kappa.Const) Route {
	return r.routes[k]
}

func (r *RouteBase) GetPath() string {
	return r.path
}

func (r *RouteBase) SetPath(path string) string {
	r.path = path
	return path
}

func (r *RouteBase) SetRouteSubstrate(sub *RouteSubstrate) {
	r.substrate = sub
}

func (r *RouteBase) GetRoute(k kappa.Const) Route {
	return r.substrate.Get(k)
}
