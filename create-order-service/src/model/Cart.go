package entities

import "github.com/supersida159/e-commerce/create-order/common"

// Cart represents a shopping cart.

type Cart struct {
	common.SQLModel `json:",inline"`
	UserID          int         `gorm:"index;column:UserID" json:"user_id"`
	User            *User       `gorm:"foreignKey:UserID" json:"user"`
	Items           []*CartItem `gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
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
