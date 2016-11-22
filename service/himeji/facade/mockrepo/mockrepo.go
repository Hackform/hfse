package mockrepo

import (
	"github.com/Hackform/hfse/service/himeji"
)

type (
	MockRepoFacade struct {
	}
)

func New() *MockRepoFacade {
	return &MockRepoFacade{}
}

func (f *MockRepoFacade) Connect(done chan<- bool) {

}

func (f *MockRepoFacade) Close() {

}

func (f *MockRepoFacade) Insert(done chan<- bool, collection string, data himeji.Data) {

}

func (f *MockRepoFacade) Query(done chan<- bool, collection string, query himeji.Bounds, result []himeji.Data) {

}

func (f *MockRepoFacade) QuerySingle(done chan<- bool, collection string, query himeji.Bounds, result *himeji.Data) {

}
