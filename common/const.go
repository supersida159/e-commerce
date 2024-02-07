package common

const (
	DbTypeRestaurant = 1
	DbTypeFood       = 2
	DbTypeCategory   = 3
	DbTypeUser       = 4
)
const (
	CurrentUser = "user"
)

type Requester interface {
	GetUserID() int
	GetRole() string
	GetEmail() string
}
