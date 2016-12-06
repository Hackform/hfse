package hfse

import (
	"fmt"
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/route"
	"github.com/Hackform/hfse/service"
	"github.com/labstack/echo"
	"time"
)

type (
	Hfse struct {
		server      *echo.Echo
		services    *service.ServiceSubstrate
		serviceList []kappa.Const
	}
)

func New() *Hfse {
	return &Hfse{
		server:   echo.New(),
		services: service.New(),
	}
}

func (h *Hfse) Start(url string) {
	for _, i := range h.serviceList {
		h.services.Get(i).Start()
	}
	fmt.Println("server starting")
	h.server.Logger.Fatal(h.server.Start(url))
}

func (h *Hfse) Shutdown() {
	for _, i := range h.serviceList {
		h.services.Get(i).Shutdown()
	}
	h.server.Shutdown(15 * time.Second)
}

func (h *Hfse) Provide(s service.Service) kappa.Const {
	s.SetServiceSubstrate(h.services)
	k := h.services.Set(s)
	h.serviceList = append(h.serviceList, k)
	return k
}

func (h *Hfse) Get(k kappa.Const) service.Service {
	return h.services.Get(k)
}

func (h *Hfse) Register(r route.Route) {
	r.SetServiceSubstrate(h.services)
	g := h.server.Group(r.GetPath(), r.Middleware()...)
	r.Register(g)
}

func (h *Hfse) Use(m ...echo.MiddlewareFunc) {
	h.server.Use(m...)
}
