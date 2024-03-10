package common

const (
	DbTypeRestaurant = 1
	DbTypeFood       = 2
	DbTypeCategory   = 3
	DbTypeUser       = 4
	DbTypeProduct    = 5
	DbTypeOrder      = 6
)
const (
	CurrentUser = "user"
)

type Requester interface {
	GetUserID() int
	GetRole() string
	GetEmail() string
	Mask(hideID bool)
}
