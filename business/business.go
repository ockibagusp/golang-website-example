package business

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// gorm.Model
type Model struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ObjectMetadata struct{}

type InternalContext struct {
	TrackerID string
}

func NewInternalContext(trackerID string) InternalContext {
	return InternalContext{
		TrackerID: trackerID,
	}
}

func (ic InternalContext) ToContext() context.Context {
	ctx := context.WithValue(context.Background(), "tracker_id", ic.TrackerID)
	return ctx
}
