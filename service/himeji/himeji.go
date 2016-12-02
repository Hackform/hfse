package himeji

import "github.com/Hackform/hfse/service"

type (
	Himeji struct {
		service.ServiceBase
		repo RepoFacade
	}

	RepoFacade interface {
		Connect(done chan<- bool)
		Close()
		Insert(done chan<- bool, collection string, data *Data)
		Query(done chan<- bool, collection string, query Bounds, result *Data)
		QueryId(done chan<- bool, collection string, query string, result *Data)
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

	Data struct {
		Value interface{}
	}

	Error string
)

func New(repo RepoFacade) *Himeji {
	return &Himeji{
		repo: repo,
	}
}

func (h *Himeji) Connect() <-chan bool {
	done := make(chan bool)
	go h.repo.Connect(done)
	return done
}

func (h *Himeji) Close() {
	h.repo.Close()
}

func (h *Himeji) Insert(collection string, data *Data) <-chan bool {
	done := make(chan bool)
	go h.repo.Insert(done, collection, data)
	return done
}

func (h *Himeji) Query(collection string, query Bounds, result *Data) <-chan bool {
	done := make(chan bool)
	go h.repo.Query(done, collection, query, result)
	return done
}

func (h *Himeji) QueryId(collection string, query string, result *Data) <-chan bool {
	done := make(chan bool)
	go h.repo.QueryId(done, collection, query, result)
	return done
}

func (e Error) Error() string {
	return string(e)
}
