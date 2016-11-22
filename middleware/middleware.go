package middleware

import (
	"github.com/labstack/echo"
)

type (
	Middleware []echo.MiddlewareFunc
)
