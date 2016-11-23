package hfse

import (
	"net/http"
	// "net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func TestHfse(t *testing.T) {
	h := New()

	h.Use(middleware.Logger())
	h.Use(middleware.Recover())

	// Start server
	go h.Start(":8080")
	time.Sleep(32 * time.Millisecond)
	h.Shutdown()
}
