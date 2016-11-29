package pionen

import (
	"github.com/Hackform/hfse/kappa"
)

type (
	Pionen struct {
		id kappa.Const
	}
)

func (h *Pionen) SetId(id kappa.Const) kappa.Const {
	h.id = id
	return h.id
}

func (h *Pionen) GetId() kappa.Const {
	return h.id
}
