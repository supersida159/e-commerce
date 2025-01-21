package usecase_user

import (
	"context"
	"fmt"

	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/pkg/tokenprovider"
	"github.com/supersida159/e-commerce/api-services/src/users/entities_user"
)

type LoginStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities_user.User, error)
}

type TokenCfg interface {
	// GetAtExp() int
	// GetRtExp() int
}
type Hasher interface {
	Hash(string) string
}

type LoginBusiness struct {
	appCtx        app_context.AppContext
	storeUser     LoginStorage
	tokenProvider tokenprovider.Provider
	hasher        Hasher
	expiry        int
}

func NewLoginBusiness(appCtx app_context.AppContext, storeUser LoginStorage, tokenProvider tokenprovider.Provider, hasher Hasher, expiry int) *LoginBusiness {
	return &LoginBusiness{
		appCtx:        appCtx,
		storeUser:     storeUser,
		tokenProvider: tokenProvider,
		hasher:        hasher,
		expiry:        expiry,
	}
}

func (b *LoginBusiness) Login(ctx context.Context, data *entities_user.UserLogin) (*entities_user.Account, *common.AppError) {
	user, err := b.storeUser.FindUser(ctx, map[string]interface{}{"email": data.Email})

	fmt.Println("data", data)

	if err != nil {
		return nil, common.ErrUserNotExist(err)
	}

	passHash := b.hasher.Hash(data.Password + user.Salt)
	if passHash != user.Password {
		return nil, common.ErrWrongPasswordOrUsername()
	}

	payload := &tokenprovider.TokenPayload{
		UserId: user.ID,
		Role:   user.Role,
	}

	accessToken, err := b.tokenProvider.Generate(payload, b.expiry)
	if err != nil {
		return nil, common.ErrInternalServerError(err)
	}
	refreshToken, err := b.tokenProvider.Generate(payload, b.expiry*2)
	if err != nil {
		return nil, common.ErrInternalServerError(err)
	}
	account := entities_user.NewAccount(accessToken, refreshToken)
	return account, nil
}
