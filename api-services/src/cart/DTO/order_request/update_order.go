package order_request

import (
	"github.com/supersida159/e-commerce/api-services/common"
)

type Cart struct {
	common.SQLModel `json:",inline"`
}
