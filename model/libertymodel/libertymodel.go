package libertymodel

import (
	"github.com/Hackform/hfse/service/himeji"
)

type (
	UserId struct {
		Id string `json:"id" bson:"_id"`
	}

	UserPassword struct {
		Password string `json:"password"`
	}

	PublicUser struct {
		UserId
		Name string `json:"name" bson:"name"`
	}

	UserPermissions struct {
		AccessLevel uint8   `json:"accesslevel" bson:"accesslevel"`
		AccessTags  []uint8 `json:"accesstags" bson:"accesstags"`
	}

	PrivateUser struct {
		UserPermissions
		Hash []byte `json:"hash" bson:"hash"`
		Salt []byte `json:"salt" bson:"salt"`
	}

	ModelUser struct {
		PublicUser
		PrivateUser
	}

	PostUser struct {
		PublicUser
		UserPassword
	}

	//////////////
	// Requests //
	//////////////

	RequestPublicUser struct {
		Value PublicUser `json:"data"`
	}

	RequestModelUser struct {
		Value ModelUser `json:"data"`
	}

	RequestPostUser struct {
		Value PostUser `json:"data"`
	}
)

const (
	collection = "Users"
)

//////////////
// Handlers //
//////////////

func GetUser(repo *himeji.Himeji, userid string, result *himeji.Data) <-chan bool {
	return repo.QueryId(collection, userid, result)
}

func StoreUser(repo *himeji.Himeji, user *himeji.Data) <-chan bool {
	return repo.Insert(collection, user)
}
