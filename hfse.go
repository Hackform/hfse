package hfse

import (
  // "net/http"
	"github.com/labstack/echo"
)

type Hfse struct {
  server *echo.Echo
}

func New() *Hfse {
	return &Hfse{
    server: echo.New(),
  }
}

func (h *Hfse) Use(m echo.MiddlewareFunc) {
  h.server.Use(m)
}

func (h *Hfse) Start(url string) {
  h.server.Logger.Fatal(h.server.Start(url))
}
