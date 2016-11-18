package hfse

import (
  "testing"
  "net/http"
	"github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
)

func TestTimeConsuming(t *testing.T) {
  h := New()

  h.Use(middleware.Logger())
  h.Use(middleware.Recover())

  h.server.GET("/", func(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!\n")
  })

  // Start server
  h.Start(":8080")
}
