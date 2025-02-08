package entities_upload

import "github.com/supersida159/e-commerce/api-services/common"

type UploadImg struct {
	common.SQLModel `json:",inline"`
	Logo            *common.Image  `json:"logo" gorm:"column:logo;"`
	Cover           *common.Images `json:"cover" gorm:"column:cover;"`
}

func (UploadImg) TableName() string {
	return "restaurants"
}

type CreateUpload struct {
	Logo  *common.Image  `json:"logo" gorm:"column:logo;"`
	Cover *common.Images `json:"cover" gorm:"column:cover;"`
}

func (CreateUpload) TableName() string {
	return UploadImg{}.TableName()
}

type UpdateUpload struct {
	Logo  *common.Image  `json:"logo" gorm:"column:logo;"`
	Cover *common.Images `json:"cover" gorm:"column:cover;"`
}
