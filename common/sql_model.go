package common

import (
	"time"
)

type SQLModel struct {
	ID        int        `json:"-" gorm:"column:id;not null;primary_key;unique"`
	FakeId    *UID       `json:"id" gorm:"-;"`
	Status    int        `json:"status" gorm:"column:status;default:1;"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"updated_at"`
}

func (m *SQLModel) GenUID(dbType int) {
	uid := NewUID(uint32(m.ID), dbType, 1)
	m.FakeId = uid
}
func (m *SQLModel) DeID() int {
	uid := m.ID
	m.FakeId = &UID{
		localID:    uint32(uid),
		objecttype: 0,
		shardID:    0,
	}
	return int(uid)
}