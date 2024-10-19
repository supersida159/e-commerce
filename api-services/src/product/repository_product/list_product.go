package repositoryproduct

import (
	"context"
	"fmt"

	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/src/product/entities_product"
	"gorm.io/gorm"
)

func (s *sqlStore) ListProduct(ctx context.Context,
	conditions map[string]interface{},
	filter *entities_product.ListProductReq,
	paging *common.Paging,
	moreInfo ...string,
) ([]entities_product.ListProductRes, error) {
	var result []entities_product.ListProductRes

	db := s.db.Table(entities_product.Product{}.TableName()).Where(conditions)

	for _, info := range moreInfo {
		db = db.Preload(info)
	}
	db = buildQuery(db, filter)

	// Find total count
	if err := db.Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	// Apply pagination conditions

	if paging.FakeCusor != "" {
		uid, err := common.FromBase58(paging.FakeCusor)
		if err == nil {
			db = db.Where("id < ?", uid.GetLocalID())
		} else {
			db = db.Offset((paging.Page - 1) * paging.Limit)
		}
	}

	// Fetch the actual data
	if err := db.
		Limit(paging.Limit).
		Order("id desc").
		Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return result, nil
}

func buildQuery(db *gorm.DB, filter *entities_product.ListProductReq) *gorm.DB {
	if filter.Name != "" {
		db = db.Where("name = ?", filter.Name)
	}
	if filter.Code != "" {
		db = db.Where("code = ?", filter.Code)
	}
	if filter.Category != "" {
		db = db.Where("category = ?", filter.Category)
	}
	if filter.Brand != "" {
		db = db.Where("brand = ?", filter.Brand)
	}
	if filter.Active != nil {
		db = db.Where("active = ?", *filter.Active)
	}
	if filter.FakeId != nil {
		fmt.Println("GetLocalID:", filter.FakeId.GetLocalID())
		db.Where("id = ?", filter.FakeId.GetLocalID())
	}
	if filter.SearchTerm != "" {
		db.Where("name LIKE ?", "%"+filter.SearchTerm+"%")
	}

	// Add conditions for CreatedAt and UpdatedAt
	if filter.CreatedAt != nil && !filter.CreatedAt.IsZero() {
		db = db.Where("created_at >= ?", filter.CreatedAt)
	}
	if filter.UpdatedAt != nil && !filter.UpdatedAt.IsZero() {
		db = db.Where("updated_at >= ?", filter.UpdatedAt)
	}

	// Add conditions for other fields as needed...

	return db
}
