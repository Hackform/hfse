package pearl

import (
	"github.com/Hackform/hfse/kappa"
)

type (
	Pearl struct {
		id kappa.Const
	}
)

func (h *Pearl) SetId(id kappa.Const) kappa.Const {
	h.id = id
	return h.id
}

func (h *Pearl) GetId() kappa.Const {
	return h.id
}
