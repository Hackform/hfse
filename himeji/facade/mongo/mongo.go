package mongo

import (
	"github.com/Hackform/hfse/himeji"
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

func (f *MongoFacade) Connect() {
	session, err := mgo.Dial(f.dbInfo.url)
	if err != nil {
		panic(err)
	}
	f.session = session
	f.database = session.DB(f.dbInfo.name)
	f.database.Login(f.dbInfo.user, f.dbInfo.pass)
}

func (f *MongoFacade) Close() {
	f.session.Close()
}

func (f *MongoFacade) Insert(collection string, query himeji.Bounds, data himeji.Data) {
	f.database.C(collection).Upsert(f.boundFormat(query), bson.M{"$set": data})
}

func (f *MongoFacade) Query(collection string, query himeji.Bounds, result []himeji.Data) {
	q := f.database.C(collection).Find(f.boundFormat(query))
	q.Iter().All(result)
}

func (f *MongoFacade) QuerySingle(collection string, query himeji.Bounds, result *himeji.Data) {
	q := f.database.C(collection).Find(f.boundFormat(query))
	q.One(result)
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
		default:
			panic("condition invalid")
		}
	}
	return m
}
