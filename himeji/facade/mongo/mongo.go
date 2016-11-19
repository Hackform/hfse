package mongo

import (
	"gopkg.in/mgo.v2"
)

type (
	MongoFacade struct {
		url     string
		session *mgo.Session
	}
)

func New(url string) *MongoFacade {
	return &MongoFacade{
		url: url,
	}
}

func (f *MongoFacade) Connect() {
	session, err := mgo.Dial(f.url)
	if err != nil {
		panic(err)
	}
	f.session = session
}
