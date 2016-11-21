package route

import (
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/middleware"
  "github.com/labstack/echo"
)

type (
	Route interface {
		Register(*Group)
    Middleware() []middleware.Middleware
	}

  Group echo.Group
)
