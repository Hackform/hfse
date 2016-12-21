package pionen

import (
	"errors"
	"fmt"
	"github.com/Hackform/hfse/kappa"
	"github.com/Hackform/hfse/model/libertymodel"
	"github.com/Hackform/hfse/service"
	"github.com/Hackform/hfse/service/himeji"
	"github.com/Hackform/hfse/service/pionen/access"
	"github.com/Hackform/hfse/service/pionen/hash"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

type (
	Pionen struct {
		service.ServiceBase
		repoService kappa.Const
		signingKey  []byte
		jwt_iss     string
		jwt_hours   int
	}

	authClaim struct {
		jwt.StandardClaims
		libertymodel.Uid
		libertymodel.PublicUser
		libertymodel.UserPermissions
	}
)

func New(signingKey, issuer string, hours int, repoService kappa.Const) *Pionen {
	return &Pionen{
		repoService: repoService,
		signingKey:  []byte(signingKey),
		jwt_iss:     issuer,
		jwt_hours:   hours,
	}
}

func (p *Pionen) Start()    {}
func (p *Pionen) Shutdown() {}

func (p *Pionen) VerifyUser(username, password string) (bool, *libertymodel.ModelUser) {
	repo := p.GetService(p.repoService).(*himeji.Himeji)

	result := new(himeji.Data)
	done := libertymodel.GetUserByUsername(repo, username, result)
	if !<-done {
		return false, nil
	}
	user := result.Value.(libertymodel.ModelUser)

	if hash.Verify(password, user.Salt, user.Hash) {
		return true, &user
	}

	return false, nil
}

func (p *Pionen) GetJWT(username, password string) (string, error) {
	if verify, user := p.VerifyUser(username, password); verify {
		claims := authClaim{
			jwt.StandardClaims{
				Issuer:    p.jwt_iss,
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Hour * time.Duration(p.jwt_hours)).Unix(),
			},
			user.Uid,
			user.PublicUser,
			user.UserPermissions,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString(p.signingKey)
	} else {
		return "", errors.New("not authorized")
	}
}

func (p *Pionen) ParseJWT(tokenString string) (*authClaim, error) {
	claims := new(authClaim)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return p.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid && claims.VerifyIssuer(p.jwt_iss, true) {
		return claims, nil
	} else {
		return nil, errors.New("jwt invalid")
	}
}

func (p *Pionen) VerifyJWTLevel(tokenString string, level uint8) bool {
	claims, err := p.ParseJWT(tokenString)

	if err != nil {
		return false
	}

	if level > access.DIVIDE {
		for _, i := range claims.AccessTags {
			if i == level {
				return true
			}
		}
	} else {
		if claims.AccessLevel <= level {
			return true
		}
	}

	return false
}

func (p *Pionen) VerifyJWTUid(tokenString string, uid string) bool {
	claims, err := p.ParseJWT(tokenString)

	if err != nil {
		return false
	}

	return claims.Uid.Id == uid
}
