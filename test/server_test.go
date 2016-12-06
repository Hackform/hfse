package server

import (
	"github.com/Hackform/hfse"
	"github.com/Hackform/hfse/route/libertyroute"
	"github.com/Hackform/hfse/route/pionenroute"
	"github.com/Hackform/hfse/service/himeji"
	"github.com/Hackform/hfse/service/himeji/facade/mockrepo"
	"github.com/Hackform/hfse/service/pionen"
	// "net/http/httptest"
	// "time"
	"github.com/labstack/echo/middleware"
	"testing"
)

func TestHfse(t *testing.T) {
	h := hfse.New()

	///////////////
	// Services  //
	///////////////

	repoFacade := mockrepo.New()
	repo := h.Provide(himeji.New(repoFacade))
	auth := h.Provide(pionen.New("signing-key", "hfse", 48, repo))

	h.Register(libertyroute.New("/users", repo))
	h.Register(pionenroute.New("/auth", auth))

	////////////////
	// Middleware //
	////////////////

	h.Use(middleware.Logger())
	h.Use(middleware.Recover())

	// Start server
	defer h.Shutdown()
	h.Start(":8080")
}
