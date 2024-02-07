package repository_user

import (
	"context"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/src/users/entities"
	"gorm.io/gorm"
)

func (s *sqlStore) FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*entities.User, error) {
	db := s.db.Begin()
	db = db.Table(entities.User{}.TableName())

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}
	var user entities.User

	if err := s.db.Where(conditions).First(&user).Error; err != nil {
		db.Rollback()
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}
	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return nil, common.ErrDB(err)
	}
	return &user, nil
}
func (s *sqlStore) CreateUser(ctx context.Context, data *entities.UserCreate) error {

	db := s.db.Begin()
	if err := db.Table(data.TableName()).Create(&data).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}
	return nil
}

func (s *sqlStore) SoftDeleteUser(ctx context.Context, id int) error {
	db := s.db

	if err := db.Table(entities.User{}.TableName()).Where("id = ?", id).Updates(map[string]interface{}{"status": 0}).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}

func (s *sqlStore) UpdateUser(ctx context.Context, data *entities.UserUpdate) error {
	db := s.db
	if err := db.Table(data.TableName()).Where("id = ?", data.ID).Updates(&data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
