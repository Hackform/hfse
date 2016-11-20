package himeji

import "github.com/Hackform/hfse/kappa"

type (
	Himeji struct {
		id   kappa.Const
		repo *RepoFacade
	}

	RepoFacade interface {
		Connect()
		Insert(collection string, data Data)
		Query(collection string, query Bounds) Data
	}

	// Bounds is a query format
	Bounds []Bound

	Bound struct {
		Condition string
		Item      string
		Value     string
	}

	Data interface {
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

func (h *Himeji) Insert(collection string, data Data) {
	(*h.repo).Insert(collection, data)
}

func (h *Himeji) Query(collection string, query Bounds) Data {
	return (*h.repo).Query(collection, query)
}

func (h *Himeji) SetId(id kappa.Const) kappa.Const {
	h.id = id
	return h.id
}

func (h *Himeji) GetId() kappa.Const {
	return h.id
}
