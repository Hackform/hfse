package himeji

import "github.com/Hackform/hfse/kappa"

type (
	Himeji struct {
		id   kappa.Const
		repo *RepoFacade
	}

	RepoFacade interface {
		Connect()
	}
)

func New(repo *RepoFacade) *Himeji {
	return &Himeji{
		id:   0,
		repo: repo,
	}
}

func (h *Himeji) Connect() {
	(*h.repo).Connect()
}

func (h *Himeji) SetId(id kappa.Const) kappa.Const {
	h.id = id
	return h.id
}

func (h *Himeji) GetId() kappa.Const {
	return h.id
}
