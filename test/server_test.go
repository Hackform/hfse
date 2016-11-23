package server

import (
	"github.com/Hackform/hfse"
	"github.com/Hackform/hfse/route/liberty"
	"github.com/Hackform/hfse/service/himeji"
	"github.com/Hackform/hfse/service/himeji/facade/mockrepo"
	// "net/http/httptest"
	// "time"
	"github.com/labstack/echo/middleware"
	"testing"
)

func TestHfse(t *testing.T) {
	h := hfse.New()
	substrate := h.GetSubstrate()

	///////////////
	// Services  //
	///////////////

	repoFacade := mockrepo.New()
	repo := himeji.New(repoFacade)
	repoId := h.Provide(repo)

	libertyRoute := liberty.New("/users", substrate, repoId)
	h.Register(libertyRoute)

	////////////////
	// Middleware //
	////////////////

	h.Use(middleware.Logger())
	h.Use(middleware.Recover())

	// Start server
	repoConnect := repo.Connect()
	defer repo.Close()
	<-repoConnect
	h.Start(":8080")
	defer h.Shutdown()
}
