package entities

import (
	"strings"

	"github.com/supersida159/e-commerce/common"
	"github.com/supersida159/e-commerce/pkg/tokenprovider"
)

type UserRole int

const (
	UserRoloUser UserRole = iota
	UserRoleManager
	UserRoleAdmin
)

func (u UserRole) String() string {
	// convert const JobState to string
	return []string{"user", "manager", "admin"}[u]
}

type User struct {
	common.SQLModel `json:",inline"`
	Email           string        `json:"email" gorm:"column:email;unique;not null;index:idx_user_email"`
	Password        string        `json:"-" gorm:"column:password;"`
	FirstName       string        `json:"first_name" gorm:"column:first_name;"`
	LastName        string        `json:"last_name" gorm:"column:last_name;"`
	Role            string        `json:"role" gorm:"column:role;"`
	Salt            string        `json:"-" gorm:"column:salt;"`
	Avatar          *common.Image `json:"avatar,omitempty" gorm:"column:avatar;type:json"`
}

func (u *User) GetUserID() int {
	return u.ID
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRole() string {
	return u.Role
}

func (User) TableName() string {
	return "users"
}
func (u *User) Mask(isAdmin bool) {
	if isAdmin {
		u.GenUID(common.DbTypeUser)
	}
}

// type RestaurantUpdate struct {
// 	Name  *string        `json:"name"`
// 	Addr  *string        `json:"addr"`
// 	Logo  *common.Image  `json:"logo" gorm:"column:logo;"`
// 	Cover *common.Images `json:"cover" gorm:"column:cover;"`
// }

// func (RestaurantUpdate) TableName() string {
// 	return Restaurant{}.TableName()
// }

type UserCreate struct {
	common.SQLModel `json:",inline"`
	Email           string         `json:"email" gorm:"column:email;"`
	Password        string         `json:"password" gorm:"column:password;"`
	FirstName       string         `json:"first_name" gorm:"column:first_name;"`
	LastName        string         `json:"last_name" gorm:"column:last_name;"`
	Role            string         `json:"-" gorm:"column:role;"`
	Salt            string         `json:"-" gorm:"column:salt;"`
	Avatar          *common.Images `json:"avatar,omitempty" gorm:"column:avatar;type:json"`
}

func (UserCreate) TableName() string {
	return User{}.TableName()
}
func (res *User) Validate() error {
	res.Email = strings.TrimSpace(res.Email)
	if len(res.Email) == 0 {
		return common.ErrNameCannotBeEmpty
	}
	return nil
}

func (u *UserCreate) Mask(isAdminOrOwner bool) {
	u.GenUID(common.DbTypeUser)
}

type UserUpdate struct {
	common.SQLModel `json:",inline"`
	Email           string        `json:"email" gorm:"column:email;"`
	Password        string        `json:"password" gorm:"column:password;"`
	FirstName       string        `json:"first_name" gorm:"column:first_name;"`
	LastName        string        `json:"last_name" gorm:"column:last_name;"`
	Role            string        `json:"-" gorm:"column:role;"`
	Salt            string        `json:"-" gorm:"column:salt;"`
	Avatar          *common.Image `json:"avatar,omitempty" gorm:"column:avatar;type:json"`
}

type UpdatePermission struct {
	UpdateEmail string `json:"updateEmail" gorm:"-"`
	RoleUpdate  string `json:"roleUpdate" gorm:"-"`
}

func (UserUpdate) TableName() string {
	return User{}.TableName()
}

type Account struct {
	AccessToken  tokenprovider.Token `json:"access_token"`
	RefreshToken tokenprovider.Token `json:"refresh_token"`
}

func NewAccount(accessToken, refreshToken *tokenprovider.Token) *Account {
	return &Account{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}
}

type UserLogin struct {
	Email    string `json:"email" gorm:"column:email;"`
	Password string `json:"password" gorm:"column:password;"`
}
