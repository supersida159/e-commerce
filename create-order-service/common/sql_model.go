package common

import (
	"fmt"
	"time"
)

type SQLModel struct {
	ID        int        `json:"real_id,omitempty" gorm:"column:id;not null;primary_key;unique"`
	FakeId    *UID       `json:"id,omitempty" gorm:"-;" validate:"-"`
	Status    int        `json:"status" gorm:"column:status;default:1;"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"deleted_at"`
}

func (m *SQLModel) GenUID(dbType int) {
	uid := NewUID(uint32(m.ID), dbType, 1)
	m.FakeId = uid
	m.ID = 0
}
func (m *SQLModel) DeID() {
	m.FakeId = nil
	fmt.Print("id:", m.ID)

}
