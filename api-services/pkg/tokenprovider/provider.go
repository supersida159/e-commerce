package tokenprovider

import (
	"time"
)

type Provider interface {
	Generate(data *TokenPayload, expiry int) (*Token, error)
	Validate(string) (*TokenPayload, error)
}

type Token struct {
	AccessToken string    `json:"access_token"`
	Created     time.Time `json:"created"`
	Expiry      int       `json:"expiry"`
}

type TokenPayload struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
}
