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
	// If FakeId is nil, there's nothing to do
	if m.FakeId == nil {
		fmt.Println("FakeId is nil, no action needed")
		return
	}

	// Get the base58 string representation from FakeId
	uidString := m.FakeId.String() // This converts FakeId (UID) to string (base58 encoded)

	// Decode the FakeId string (base58) back to a UID
	uid, err := FromBase58(uidString)
	if err != nil {
		// Handle the error if the UID decoding fails
		fmt.Println("Failed to decode FakeId:", err)
		return
	}

	// Extract the original fields (localID, objecttype, shardID)
	m.ID = int(uid.GetLocalID()) // Assuming you want to restore the ID field with localID
	m.FakeId = nil               // Reset the FakeId after extraction

	// Optionally, log or handle the extracted values (localID, objecttype, shardID)
	fmt.Printf("Decomposed UID - ID: %d, ObjectType: %d, ShardID: %d\n", m.ID, uid.GetObjectType(), uid.GetShardID())
}
