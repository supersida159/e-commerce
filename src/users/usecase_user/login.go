package usecase_user

import (
	"context"
	"fmt"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/pkg/tokenprovider"
	"github.com/supersida159/e-commerce/src/users/entities"
)

type LoginStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities.User, error)
}

type TokenCfg interface {
	// GetAtExp() int
	// GetRtExp() int
}
type Hasher interface {
	Hash(string) string
}

type LoginBusiness struct {
	appCtx        app_context.Appcontext
	storeUser     LoginStorage
	tokenProvider tokenprovider.Provider
	hasher        Hasher
	expiry        int
}

func NewLoginBusiness(appCtx app_context.Appcontext, storeUser LoginStorage, tokenProvider tokenprovider.Provider, hasher Hasher, expiry int) *LoginBusiness {
	return &LoginBusiness{
		appCtx:        appCtx,
		storeUser:     storeUser,
		tokenProvider: tokenProvider,
		hasher:        hasher,
		expiry:        expiry,
	}
}

func (b *LoginBusiness) Login(ctx context.Context, data *entities.UserLogin) (*entities.Account, error) {
	user, err := b.storeUser.FindUser(ctx, map[string]interface{}{"email": data.Email})

	fmt.Println("data", data)

	if err != nil {
		return nil, common.ErrCannotGetEntity(entities.EntityName, err)
	}

	passHash := b.hasher.Hash(data.Password + user.Salt)
	if passHash != user.Password {
		return nil, &common.ErrUserNameOrPasswordInvalid
	}

	payload := &tokenprovider.TokenPayload{
		UserId: user.ID,
		Role:   user.Role,
	}

	accessToken, err := b.tokenProvider.Generate(payload, b.expiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}
	refreshToken, err := b.tokenProvider.Generate(payload, b.expiry*2)
	if err != nil {
		return nil, common.ErrInternal(err)
	}
	account := entities.NewAccount(accessToken, refreshToken)
	return account, nil
}
