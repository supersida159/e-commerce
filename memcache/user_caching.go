package memcache

import (
	"context"
	"fmt"

	"github.com/supersida159/e-commerce/src/users/entities_user"
)

func (uc *memcaching) FindUser(ctx context.Context, condition map[string]interface{}, moreInfores ...string) (*entities_user.User, error) {
	userID := condition["id"].(int)

	userInCache, err := uc.store.Read(fmt.Sprintf("user-%d", userID))
	if err != nil {
		return userInCache.(*entities_user.User), nil
	}
	userInRealStore, err := uc.realStore.FindUser(ctx, condition, moreInfores...)
	if err != nil {
		return nil, err
	}
	go func() {
		uc.store.Write(fmt.Sprintf("user-%d", userID), userInRealStore)
	}()
	return userInRealStore, nil
}
