package usecase_user

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/app_context"
	"github.com/supersida159/e-commerce/src/users/entities"
)

type UpdatePermissionStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities.User, error)
	UpdateUserAddmin(ctx context.Context, data *entities.UserUpdate) error
}

type updatePermissionBusiness struct {
	storeUser UpdatePermissionStorage
	appCtx    app_context.Appcontext
	hasher    Hasher
}

func NewUpdatePermissionBusiness(appCtx app_context.Appcontext, storeUser UpdatePermissionStorage, hasher Hasher) *updatePermissionBusiness {
	return &updatePermissionBusiness{
		appCtx:    appCtx,
		storeUser: storeUser,
		hasher:    hasher,
	}
}

func (b *updatePermissionBusiness) UpdateUserPermission(ctx context.Context, data *entities.UserUpdate) error {
	if err := b.storeUser.UpdateUserAddmin(ctx, data); err != nil {
		return common.ErrDB(err)
	}
	return nil
}
