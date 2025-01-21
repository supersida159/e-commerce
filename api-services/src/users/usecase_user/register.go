package usecase_user

import (
	"context"

	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/src/users/entities_user"
)

type RegisterStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities_user.User, error)
	CreateUser(ctx context.Context, data *entities_user.UserCreate) error
}

type RegisterBusiness struct {
	storeUser RegisterStorage
	appCtx    app_context.AppContext
	hasher    Hasher
}

func NewRegisterBusiness(appCtx app_context.AppContext, storeUser RegisterStorage, hasher Hasher) *RegisterBusiness {
	return &RegisterBusiness{
		appCtx:    appCtx,
		storeUser: storeUser,
		hasher:    hasher,
	}
}

func (b *RegisterBusiness) Register(ctx context.Context, data *entities_user.UserCreate) error {

	user, err := b.storeUser.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if user != nil {
		return common.ErrUserNameExists(err)
	}

	salt := common.GenSalt(50)

	data.Password = b.hasher.Hash(data.Password + salt)
	data.Salt = salt
	data.Role = "user"
	data.Status = 1
	if err := b.storeUser.CreateUser(ctx, data); err != nil {
		return common.ErrDB(err)
	}
	return nil
}
