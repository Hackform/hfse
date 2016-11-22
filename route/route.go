package route

import (
	"github.com/labstack/echo"
)

type (
	Route interface {
		Register(*echo.Group)
		Middleware() []echo.MiddlewareFunc
	}
)
