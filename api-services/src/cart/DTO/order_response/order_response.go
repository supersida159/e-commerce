package order_response

import (
	"github.com/supersida159/e-commerce/api-services/common"
	"github.com/supersida159/e-commerce/api-services/src/product/entities_product"
	"github.com/supersida159/e-commerce/api-services/src/users/entities_user"
)

type CartResponse struct {
	common.SQLModel `json:",inline"`
	User            *entities_user.User          `gorm:"foreignKey:UserID" json:"user"`
	Items           []*entities_product.CartItem `gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
}
