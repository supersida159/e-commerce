package usecase_user

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/src/users/entities"
)

type UpdateStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities.User, error)
	UpdateUser(ctx context.Context, data *entities.UserUpdate) error
}

type UpdateBusiness struct {
	storeUser UpdateStorage
	appCtx    app_context.Appcontext
	hasher    Hasher
}

func NewUpdateBusiness(appCtx app_context.Appcontext, storeUser UpdateStorage, hasher Hasher) *UpdateBusiness {
	return &UpdateBusiness{
		appCtx:    appCtx,
		storeUser: storeUser,
		hasher:    hasher,
	}
}

func (b *UpdateBusiness) UpdateUser(ctx context.Context, data *entities.UserUpdate) error {
	if data.Password != "" {
		salt := common.GenSalt(50)
		data.Password = b.hasher.Hash(data.Password + salt)
		data.Salt = salt
	}
	if err := b.storeUser.UpdateUser(ctx, data); err != nil {
		return common.ErrDB(err)
	}
	return nil
}
