package route

import (
	"github.com/Hackform/hfse/kappa"
)

type (
	Route interface {
		SetId(id string) string
		GetId() string
	}
)
