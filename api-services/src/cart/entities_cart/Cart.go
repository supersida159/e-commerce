package entities_carts

import (
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/src/product/entities_product"
	"github.com/supersida159/e-commerce/api-services/src/users/entities_user"
)

// Cart represents a shopping cart.

const EntityName = "Cart"

type Cart struct {
	common.SQLModel `json:",inline"`
	UserID          int                          `gorm:"index;column:UserID" json:"user_id"`
	User            *entities_user.User          `gorm:"foreignKey:UserID" json:"user"`
	Items           []*entities_product.CartItem `gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
}

// CartItem represents an item in the shopping cart.

func (c Cart) TableName() string {
	return "Cart"
}

func (o *Cart) Mask(hideID bool) {
	if hideID {
		o.GenUID(common.DbTypeOrder)
	} else if !hideID {
		o.DeID()
	}
}
