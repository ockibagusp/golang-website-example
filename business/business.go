package business

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

// gorm.Model
type Model struct {
	ID uint `gorm:"primarykey"`
	// rittneje: `time.Time` to `sql.Scanner`
	// -> https://github.com/mattn/go-sqlite3/issues/951
	CreatedAt sql.Scanner    `json:"-"`
	UpdatedAt sql.Scanner    `json:"-"`
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
