package himeji

import (
	"github.com/Hackform/hfse/kappa"
)

type (
	Himeji struct {
		id kappa.Const
	}
)

func New() *Himeji {
	return &Himeji{
		id: 0,
	}
}

func (h *Himeji) SetId(id kappa.Const) kappa.Const {
	h.id = id
	return h.id
}

func (h *Himeji) GetId() kappa.Const {
	return h.id
}
