package location

import "golang-website-example/business"

type (
	Location struct {
		ID       uint   `form:"id"`
		Location string `form:"location"`

		business.ObjectMetadata
	}
)

// TableName name: string
func (Location) TableName() string {
	return "locations"
}
