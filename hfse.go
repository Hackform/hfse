package hfse

import (
  // "net/http"
  "time"
	"github.com/labstack/echo"
  "github.com/Hackform/hfse/kappa"
  "github.com/Hackform/hfse/service"
  "github.com/Hackform/hfse/middleware"
)

type (
  Hfse struct {
    server *echo.Echo
    serviceKappa *kappa.Kappa
    services map[kappa.Const]*service.Service
  }
)

func New() *Hfse {
	return &Hfse{
    server: echo.New(),
    serviceKappa: kappa.New(),
    services: make(map[kappa.Const]*service.Service),
  }
}

func (h *Hfse) Start(url string) {
  h.server.Logger.Fatal(h.server.Start(url))
}

func (h *Hfse) Shutdown() {
  h.server.Shutdown(15 * time.Second)
}


func (h *Hfse) Register(s *service.Service) kappa.Const {
  k := h.serviceKappa.Get()
  (*s).SetId(k)
  h.services[k] = s
  return k
}

func (h *Hfse) Use(m middleware.Middleware) {
  h.server.Use(m.(echo.MiddlewareFunc))
}
