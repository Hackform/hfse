package route

import (
	"github.com/Hackform/hfse/kappa"
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
		Register(*echo.Group)
		Middleware() []echo.MiddlewareFunc
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
