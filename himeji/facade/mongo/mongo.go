package mongo

import (
	"github.com/Hackform/hfse/himeji"
	"gopkg.in/mgo.v2"
)

type (
	mongoDbInfo struct {
		url  string
		name string
		user string
		pass string
	}

	MongoFacade struct {
		dbInfo   mongoDbInfo
		database *mgo.Database
	}
)

func New(url, name, user, pass string) *MongoFacade {
	return &MongoFacade{
		dbInfo: mongoDbInfo{
			url:  url,
			name: name,
			user: user,
			pass: pass,
		},
	}
}

func (f *MongoFacade) Connect() {
	session, err := mgo.Dial(f.dbInfo.url)
	if err != nil {
		panic(err)
	}
	f.database = session.DB(f.dbInfo.name)
	f.database.Login(f.dbInfo.user, f.dbInfo.pass)
}

func (f *MongoFacade) Insert(collection string, data himeji.Data) {
	f.database.C(collection).Insert(data)
}

func (f *MongoFacade) Query(collection string, query himeji.Bounds) himeji.Data {
	for _, bound := range query {

	}
	return f.database.C(collection).Find()
}

func (f *MongoFacade) boundFormat(bound himeji.Bound) (string, interface{}) {
	condition := bound.Condition
	item := bound.Item
	value := bound.Value
}
