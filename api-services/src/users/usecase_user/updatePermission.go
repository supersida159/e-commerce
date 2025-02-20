package usecase_user

import (
	"context"

	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/src/users/entities_user"
)

type UpdatePermissionStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities_user.User, *common.AppError)
	UpdateUser(ctx context.Context, data *entities_user.UserUpdate) *common.AppError
}

type updatePermissionBusiness struct {
	storeUser UpdatePermissionStorage
	appCtx    app_context.AppContext
	hasher    Hasher
}

func NewUpdatePermissionBusiness(appCtx app_context.AppContext, storeUser UpdatePermissionStorage, hasher Hasher) *updatePermissionBusiness {
	return &updatePermissionBusiness{
		appCtx:    appCtx,
		storeUser: storeUser,
		hasher:    hasher,
	}
}

func (b *updatePermissionBusiness) UpdateUserPermission(ctx context.Context, data *entities_user.UserUpdate) *common.AppError {
	if err := b.storeUser.UpdateUser(ctx, data); err != nil {
		return common.ErrDB(err)
	}
	return nil
}
