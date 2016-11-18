package hfse

import (
  "testing"
  "time"
  "net/http"
	"github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
)

func TestHfse(t *testing.T) {
  h := New()

  h.Use(middleware.Logger())
  h.Use(middleware.Recover())

  h.server.GET("/", func(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!\n")
  })

  // Start server
  go h.Start(":8080")
  time.Sleep(250 * time.Millisecond)
  h.Shutdown()
}
