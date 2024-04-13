package entities_product

import "github.com/supersida159/e-commerce/common"

const (
	EntityName = "Product"
)

type Product struct {
	common.SQLModel `json:",inline"`
	Description     string        `json:"description" gorm:"column:description"`
	Name            string        `json:"name" gorm:"Index:idx_product_name(255),not null;column:name;type:varchar(255)"`
	Code            string        `json:"code" gorm:"uniqueIndex:idx_product_code(255),not null;column:code;type:varchar(255)"`
	Price           string        `json:"price" gorm:"column:price;not null"`
	Active          *bool         `json:"active" gorm:"default:true"`
	Quantity        int           `json:"quantity" gorm:"column:quantity;not null"`
	ProductImages   *ProductImage `json:"images" gorm:"column:product_image;type:json"`
	Category        string        `json:"category" gorm:"column:category;type:varchar(255)"`
	Brand           string        `json:"brand" gorm:"column:brand;type:varchar(255)"`
	Sale            int           `json:"sale" gorm:"column:sale;not null;default:0"`
}

func (p Product) TableName() string {
	return "products"
}
func (p *Product) GetProductID() int {
	return p.ID

}

func (p *Product) GetProductName() string {
	return p.Name
}

func (p *Product) GetProductCode() string {
	return p.Code
}

func (p *Product) GetPrice() string {
	return p.Price
}

func (p *Product) GetQuantity() int {
	return p.Quantity
}

func (p *Product) GetActive() *bool {
	return p.Active
}

func (p *Product) GetDescription() string {
	return p.Description
}

// func (p *Product) GetProductImage() *ProductImage {
// 	return p.ProductImage
// }

func (p *Product) Mask(hideID bool) {
	if hideID {
		p.GenUID(common.DbTypeProduct)
	} else if !hideID {
		p.DeID()
	}
}

type ListProductReq struct {
	common.SQLModel `json:",inline"`
	Name            string `json:"name,omitempty" form:"name"`
	Code            string `json:"code,omitempty" form:"code"`
	Category        string `json:"category,omitempty" form:"category"`
	Brand           string `json:"brand,omitempty" form:"brand"`
	Active          *bool  `json:"active,omitempty" form:"active"`
	Limit           int64  `json:"-" form:"limit"`
	OrderBy         string `json:"-" form:"order_by"`
	OrderDesc       bool   `json:"-" form:"order_desc"`
}

func (p *ListProductReq) Mask(hideID bool) {
	if hideID {
		p.GenUID(common.DbTypeProduct)
	} else if !hideID {
		p.DeID()
	}
}

type ListProductRes struct {
	common.SQLModel `json:",inline"`
	Name            string        `json:"name,omitempty" form:"name"`
	Code            string        `json:"code,omitempty" form:"code"`
	Category        string        `json:"category,omitempty" form:"category"`
	ProductImages   *ProductImage `json:"images" gorm:"column:product_image;type:json"`
	Active          *bool         `json:"active" gorm:"default:true"`
	Price           string        `json:"price" gorm:"not null"`
	Brand           string        `json:"brand,omitempty" form:"brand"`
	Quantity        int           `json:"quantity" gorm:"column:quantity;not null"`
	Sale            int           `json:"sale" gorm:"column:sale;not null"`
}

func (p *ListProductRes) Mask(hideID bool) {
	if hideID {
		p.GenUID(common.DbTypeProduct)
	} else if !hideID {
		p.DeID()
	}
}
