package usecase_user

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/src/users/entities_user"
)

type RegisterStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities_user.User, error)
	CreateUser(ctx context.Context, data *entities_user.UserCreate) error
}

type RegisterBusiness struct {
	storeUser RegisterStorage
	appCtx    app_context.Appcontext
	hasher    Hasher
}

func NewRegisterBusiness(appCtx app_context.Appcontext, storeUser RegisterStorage, hasher Hasher) *RegisterBusiness {
	return &RegisterBusiness{
		appCtx:    appCtx,
		storeUser: storeUser,
		hasher:    hasher,
	}
}

func (b *RegisterBusiness) Register(ctx context.Context, data *entities_user.UserCreate) error {

	user, err := b.storeUser.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if user != nil {
		return common.ErrEntityExisted(entities_user.UserRoloUser.String(), err)
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
