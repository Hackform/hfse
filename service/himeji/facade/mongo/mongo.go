package mongo

import (
	"github.com/Hackform/hfse/service/himeji"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
		session  *mgo.Session
		database *mgo.Database
	}
)

func New(url, name, user, pass string, limit int) *MongoFacade {
	return &MongoFacade{
		dbInfo: mongoDbInfo{
			url:  url,
			name: name,
			user: user,
			pass: pass,
		},
	}
}

func (f *MongoFacade) Connect(done chan<- bool) {
	session, err := mgo.Dial(f.dbInfo.url)
	if err != nil {
		done <- false
	}
	f.session = session
	f.database = session.DB(f.dbInfo.name)
	err = f.database.Login(f.dbInfo.user, f.dbInfo.pass)
	if err != nil {
		done <- false
	}
	done <- true
}

func (f *MongoFacade) Close() {
	f.session.Close()
}

func (f *MongoFacade) Insert(done chan<- bool, collection string, query himeji.Bounds, data *himeji.Data) {
	_, err := f.database.C(collection).Upsert(f.boundFormat(query), bson.M{"$set": data.Value})
	if err != nil {
		done <- false
	}
	done <- true
}

func (f *MongoFacade) Query(done chan<- bool, collection string, query himeji.Bounds, result *himeji.Data) {
	q := f.database.C(collection).Find(f.boundFormat(query))
	err := q.Iter().All(result.Value)
	if err != nil {
		done <- false
	}
	done <- true
}

func (f *MongoFacade) QueryId(done chan<- bool, collection string, query string, result *himeji.Data) {
	q := f.database.C(collection).Find(bson.M{"_id": query})
	err := q.One(result.Value)
	if err != nil {
		done <- false
	}
	done <- true
}

func (f *MongoFacade) boundFormat(bounds himeji.Bounds) map[string]interface{} {
	m := make(map[string]interface{})
	for _, bound := range bounds {
		condition := bound.Condition
		item := bound.Item
		value := bound.Value
		switch condition {
		case "equal":
			m[item] = value
		}
	}
	return m
}
