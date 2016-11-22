package route

import (
	"github.com/labstack/echo"
)

type (
	Route interface {
		GetPath() string
		Register(*echo.Group)
		Middleware() []echo.MiddlewareFunc
	}
)
