package mongo

import (
	"time"
)

// categoryAttributeEntity represents the MongoDB document structure for category-attribute assignments
type categoryAttributeEntity struct {
	ID          string    `bson:"_id"`
	Version     int       `bson:"version"`
	CategoryID  string    `bson:"categoryId"`
	AttributeID string    `bson:"attributeId"`
	Required    bool      `bson:"required"`
	SortOrder   int       `bson:"sortOrder"`
	Filterable  *bool     `bson:"filterable,omitempty"`
	Searchable  *bool     `bson:"searchable,omitempty"`
	Enabled     bool      `bson:"enabled"`
	CreatedAt   time.Time `bson:"createdAt"`
	ModifiedAt  time.Time `bson:"modifiedAt"`
}
