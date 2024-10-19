package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/supersida159/e-commerce/api-services/pkg/tokenprovider"
)

type jwtProvider struct {
	secret string
}

func NewJwtProvider(secret string) *jwtProvider {
	return &jwtProvider{
		secret: secret,
	}
}

type myClaim struct {
	Payload tokenprovider.TokenPayload `json:"payload"`
	jwt.StandardClaims
}

func (j *jwtProvider) Generate(data *tokenprovider.TokenPayload, expiry int) (*tokenprovider.Token, error) {
	claims := myClaim{
		*data,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Second * time.Duration(expiry)).Unix(),
			IssuedAt:  time.Now().Local().Unix(),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	mytoken, err := t.SignedString([]byte(j.secret))
	if err != nil {
		return nil, err
	}
	return &tokenprovider.Token{
		AccessToken: mytoken,
		Created:     time.Now(),
		Expiry:      expiry,
	}, nil
}

func (j *jwtProvider) Validate(myToken string) (*tokenprovider.TokenPayload, error) {
	t, err := jwt.ParseWithClaims(myToken, &myClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claim, ok := t.Claims.(*myClaim)
	if !ok {
		return nil, errors.New("the token is invalid")
	}

	return &claim.Payload, nil
}
