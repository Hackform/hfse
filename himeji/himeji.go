package himeji

import "github.com/Hackform/hfse/kappa"

type (
	Himeji struct {
		id   kappa.Const
		repo *RepoFacade
	}

	RepoFacade interface {
		Connect()
		Close()
		Insert(done chan<- bool, collection string, data Data)
		Query(done chan<- bool, collection string, query Bounds, result []Data)
		QuerySingle(done chan<- bool, collection string, query Bounds, result *Data)
	}

	Bounds []Bound

	// Bound is a query format
	//
	// conditions:
	// equal
	// limit
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

func (h *Himeji) Close() {
	(*h.repo).Close()
}

func (h *Himeji) Insert(collection string, data Data) <-chan bool {
	done := make(chan bool)
	go func() { (*h.repo).Insert(done, collection, data) }()
	return done
}

func (h *Himeji) Query(collection string, query Bounds, result []Data) <-chan bool {
	done := make(chan bool)
	go func() { (*h.repo).Query(done, collection, query, result) }()
	return done
}

func (h *Himeji) QuerySingle(collection string, query Bounds, result *Data) <-chan bool {
	done := make(chan bool)
	go func() { (*h.repo).QuerySingle(done, collection, query, result) }()
	return done
}

func (h *Himeji) SetId(id kappa.Const) kappa.Const {
	h.id = id
	return h.id
}

func (h *Himeji) GetId() kappa.Const {
	return h.id
}
