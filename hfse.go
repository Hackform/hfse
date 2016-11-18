package hfse

import (
  // "net/http"
  "time"
	"github.com/labstack/echo"
)

type (
  Hfse struct {
    server *echo.Echo
    serviceKappa *Kappa
    services map[uint16]*Service
  }

  Service interface {

  }
)



func New() *Hfse {
	return &Hfse{
    server: echo.New(),
    serviceKappa: NewKappa(),
  }
}

func (h *Hfse) Start(url string) {
  h.server.Logger.Fatal(h.server.Start(url))
}

func (h *Hfse) Shutdown() {
  h.server.Shutdown(15 * time.Second)
}


func (h *Hfse) Register(s Service) {

}

func (h *Hfse) Use(m echo.MiddlewareFunc) {
  h.server.Use(m)
}
