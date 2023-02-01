package business

import (
	"context"
	"time"
)

type ObjectMetadata struct {
	CreatedAt  time.Time `gorm:"created_at"`
	CreatedBy  int       `gorm:"created_by"`
	ModifiedAt time.Time `gorm:"modified_at"`
	ModifiedBy int       `gorm:"modified_by"`
	DeletedAt  time.Time `gorm:"deleted_at"`
	DeletedBy  int       `gorm:"deleted_by"`
	Version    int       `gorm:"version"`
}

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
