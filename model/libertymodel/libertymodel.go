package libertymodel

import (
	"github.com/Hackform/hfse/service/himeji"
	"github.com/labstack/echo"
)

type (
	Uid struct {
		Id string `json:"id" bson:"_id"`
	}

	UserId struct {
		Username string `json:"username" bson:"username"`
	}

	UserPassword struct {
		Password string `json:"password"`
	}

	PublicUser struct {
		UserId
		Name string `json:"name" bson:"name"`
	}

	UserSecurity struct {
		Hash []byte `json:"hash" bson:"hash"`
		Salt []byte `json:"salt" bson:"salt"`
	}

	UserPermissions struct {
		AccessLevel uint8   `json:"accesslevel" bson:"accesslevel"`
		AccessTags  []uint8 `json:"accesstags" bson:"accesstags"`
	}

	PrivateUser struct {
		Email string `json:"email" bson:"email"`
	}

	UserInfo struct {
		PublicUser
		PrivateUser
	}

	ModelUser struct {
		Uid
		UserInfo
		UserSecurity
		UserPermissions
	}

	PostUser struct {
		UserInfo
		UserPassword
	}

	//////////////
	// Requests //
	//////////////

	RequestPublicUser struct {
		Value PublicUser `json:"data"`
	}

	RequestUserInfo struct {
		Value UserInfo `json:"data"`
	}

	RequestPostUser struct {
		Value PostUser `json:"data"`
	}

	RequestUserPassword struct {
		Value UserPassword `json:"data"`
	}
)

func GetRequestPostUser(c echo.Context) *RequestPostUser {
	user := new(RequestPostUser)
	c.Bind(user)
	return user
}

func GetRequestUserInfo(c echo.Context) *RequestUserInfo {
	user := new(RequestUserInfo)
	c.Bind(user)
	return user
}

func GetRequestUserPassword(c echo.Context) *RequestUserPassword {
	user := new(RequestUserPassword)
	c.Bind(user)
	return user
}

const (
	collection = "Users"
)

//////////////
// Handlers //
//////////////

func GetUser(repo *himeji.Himeji, uid string, result *himeji.Data) <-chan bool {
	return repo.QueryId(collection, uid, result)
}

func GetUserByUsername(repo *himeji.Himeji, username string, result *himeji.Data) <-chan bool {
	return repo.Query(collection, himeji.Bounds{
		himeji.Bound{
			Condition: "equal",
			Item:      "username",
			Value:     username,
		},
	}, result)
}

func StoreUser(repo *himeji.Himeji, user *himeji.Data) <-chan bool {
	return repo.Insert(collection, user)
}
