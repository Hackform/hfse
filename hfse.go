package hfse

import (
	// "net/http"
	"fmt"
	"time"

	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/route"
	"github.com/Hackform/hfse/service"
	"github.com/labstack/echo"
)

type (
	Hfse struct {
		server   *echo.Echo
		services *service.ServiceSubstrate
		routes   *route.RouteSubstrate
	}
)

func New() *Hfse {
	return &Hfse{
		server:   echo.New(),
		services: service.New(),
		routes:   route.New(),
	}
}

func (h *Hfse) Start(url string) {
	fmt.Println("server starting")
	h.server.Logger.Fatal(h.server.Start(url))
}

func (h *Hfse) Shutdown() {
	h.server.Shutdown(15 * time.Second)
}

func (h *Hfse) Provide(s service.Service) kappa.Const {
	s.SetServiceSubstrate(h.services)
	return h.services.Set(s)
}

func (h *Hfse) Register(r route.Route) kappa.Const {
	g := h.server.Group(r.GetPath(), r.Middleware()...)
	r.SetServiceSubstrate(h.services)
	r.SetRouteSubstrate(h.routes)
	r.Register(g)
	return h.routes.Set(r)
}

func (h *Hfse) Use(m ...echo.MiddlewareFunc) {
	h.server.Use(m...)
}
