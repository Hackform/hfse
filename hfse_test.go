package hfse

import (
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func TestHfse(t *testing.T) {
	h := New()

	h.server.Use(middleware.Logger())
	h.server.Use(middleware.Recover())

	h.server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	// Start server
	go h.Start(":8080")
	time.Sleep(32 * time.Millisecond)
	h.Shutdown()
}
