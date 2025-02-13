package usecase_oauth

import (
	"github.com/supersida159/e-commerce/api-services/pkg/app_context"
	"github.com/supersida159/e-commerce/api-services/src/users/usecase_user"
)

type OAuthBusiness struct {
	appCtx            app_context.AppContext
	StoreUser         *usecase_user.LoginBusiness
	StoreRegisterUser *usecase_user.RegisterBusiness
}

func NewOAuthBusiness(appCtx app_context.AppContext, storeUser *usecase_user.LoginBusiness, storeRegisterUser *usecase_user.RegisterBusiness) *OAuthBusiness {
	return &OAuthBusiness{
		appCtx:            appCtx,
		StoreUser:         storeUser,
		StoreRegisterUser: storeRegisterUser,
	}
}
