package mongo

import (
	"time"
)

// attributeEntity represents the MongoDB document structure
type attributeEntity struct {
	ID                string    `bson:"_id"`
	Version           int       `bson:"version"`
	Name              string    `bson:"name"`
	Slug              string    `bson:"slug"`
	Type              string    `bson:"type"`
	Unit              *string   `bson:"unit,omitempty"`
	DefaultFilterable bool      `bson:"defaultFilterable"`
	DefaultSearchable bool      `bson:"defaultSearchable"`
	SortOrder         int       `bson:"sortOrder"`
	Enabled           bool      `bson:"enabled"`
	CreatedAt         time.Time `bson:"createdAt"`
	ModifiedAt        time.Time `bson:"modifiedAt"`
}
