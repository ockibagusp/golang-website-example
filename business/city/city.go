package city

import "github.com/ockibagusp/golang-website-example/business"

type (
	City struct {
		ID   uint   `form:"id"`
		City string `form:"city"`

		business.ObjectMetadata
	}
)

// TableName name: string
func (City) TableName() string {
	return "cities"
}
