package mockrepo

import (
	"encoding/json"
	// "fmt"
	"github.com/Hackform/hfse/service/himeji"
)

var (
	db map[string]map[string]interface{}
)

type (
	MockRepoFacade struct {
	}
)

func New() *MockRepoFacade {
	return &MockRepoFacade{}
}

func (f *MockRepoFacade) Connect(done chan<- bool) {
	db = make(map[string]map[string]interface{})
	done <- true
}

func (f *MockRepoFacade) Close() {

}

func (f *MockRepoFacade) Insert(done chan<- bool, collection string, data *himeji.Data) {
	id := extractId(data.Value)
	if id == "" {
		done <- false
	} else {
		if _, ok := db[collection]; !ok {
			db[collection] = make(map[string]interface{})
		}
		db[collection][id] = data.Value
		done <- true
	}
}

func (f *MockRepoFacade) Query(done chan<- bool, collection string, query himeji.Bounds, result *himeji.Data) {
	done <- false
}

func (f *MockRepoFacade) QueryId(done chan<- bool, collection string, query string, result *himeji.Data) {
	if val, ok := db[collection]; ok {
		if val, ok := val[query]; ok {
			result.Value = val
			done <- true
		}
	}
	done <- false
}

func extractId(data interface{}) string {
	marshaled, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	dat := make(map[string]interface{})
	json.Unmarshal(marshaled, &dat)
	return dat["id"].(string)
}
