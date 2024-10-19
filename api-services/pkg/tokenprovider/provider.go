package tokenprovider

import (
	"errors"
	"time"

	"github.com/supersida159/e-commerce/api-services/common"
)

type Provider interface {
	Generate(data *TokenPayload, expiry int) (*Token, error)
	Validate(string) (*TokenPayload, error)
}

var (
	ErrNotFoundToken = common.NewCustomError(
		errors.New("not found token"),
		"not found token",
		"ErrNotFoundToken",
	)
	ErrEncodingToken = common.NewCustomError(
		errors.New("error encoding token"),
		"error encoding token",
		"ErrEncodingToken",
	)
	ErrInvalidToken = common.NewCustomError(
		errors.New("invalid token"),
		"invalid token",
		"ErrInvalidToken",
	)
)

type Token struct {
	AccessToken string    `json:"access_token"`
	Created     time.Time `json:"created"`
	Expiry      int       `json:"expiry"`
}

type TokenPayload struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
}
