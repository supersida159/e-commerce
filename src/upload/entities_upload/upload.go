package entities_Upload

import "github.com/supersida159/e-commerce/common"

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

func (UpdateUpload) TableName() string {
	return UploadImg{}.TableName()
}
func ErrCannotSaveFile(err error) *common.AppError {
	return common.NewCustomError(
		err,
		"cannot save file",
		"ErrCannotSaveFile")
}

func ErrCannotSaveImgOnDB(err error) *common.AppError {
	return common.NewCustomError(
		err,
		"cannot upload to db Image",
		"ErrCannotSaveImage")
}
