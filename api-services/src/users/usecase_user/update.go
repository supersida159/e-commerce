package usecase_user

import (
	"context"

	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/src/users/entities_user"
)

type UpdateStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities_user.User, error)
	UpdateUser(ctx context.Context, data *entities_user.UserUpdate) error
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

func (b *UpdateBusiness) UpdateUser(ctx context.Context, data *entities_user.UserUpdate) error {
	var oldData *entities_user.User

	oldData, err := b.storeUser.FindUser(ctx, map[string]interface{}{"id": data.ID})
	if err != nil {
		return common.ErrCannotGetEntity(entities_user.UserRoloUser.String(), err)
	}
	if data.NewPassword != "" {
		if data.Password != "" {
			if data.NewPassword == data.Password {
				return common.NewErrInvalidPassword()
			}

			data.Password = b.hasher.Hash(data.Password + oldData.Salt)
			if data.Password != oldData.Password {
				return common.NewErrInvalidPassword()
			}
			salt := common.GenSalt(50)

			data.Password = b.hasher.Hash(data.NewPassword + salt)

			data.Salt = salt

		} else {
			return common.NewErrInvalidPassword()
		}
	} else {
		data.Password = ""
	}
	if err := b.storeUser.UpdateUser(ctx, data); err != nil {
		return common.ErrDB(err)
	}
	return nil
}
