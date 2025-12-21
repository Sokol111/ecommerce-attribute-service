package categoryattribute

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// CategoryAttribute represents an assignment of an attribute to a category
type CategoryAttribute struct {
	ID          string
	Version     int
	CategoryID  string
	AttributeID string
	Required    bool
	SortOrder   int
	Filterable  *bool // nil means use attribute default
	Searchable  *bool // nil means use attribute default
	Enabled     bool
	CreatedAt   time.Time
	ModifiedAt  time.Time
}

// NewCategoryAttribute creates a new category-attribute assignment with validation
func NewCategoryAttribute(
	id string,
	categoryID string,
	attributeID string,
	required bool,
	sortOrder int,
	filterable *bool,
	searchable *bool,
	enabled bool,
) (*CategoryAttribute, error) {
	if err := validateCategoryAttributeData(categoryID, attributeID, sortOrder); err != nil {
		return nil, err
	}

	if id == "" {
		id = uuid.New().String()
	}

	now := time.Now().UTC()
	return &CategoryAttribute{
		ID:          id,
		Version:     1,
		CategoryID:  categoryID,
		AttributeID: attributeID,
		Required:    required,
		SortOrder:   sortOrder,
		Filterable:  filterable,
		Searchable:  searchable,
		Enabled:     enabled,
		CreatedAt:   now,
		ModifiedAt:  now,
	}, nil
}

// Reconstruct rebuilds a category attribute from persistence (no validation)
func Reconstruct(
	id string,
	version int,
	categoryID string,
	attributeID string,
	required bool,
	sortOrder int,
	filterable *bool,
	searchable *bool,
	enabled bool,
	createdAt time.Time,
	modifiedAt time.Time,
) *CategoryAttribute {
	return &CategoryAttribute{
		ID:          id,
		Version:     version,
		CategoryID:  categoryID,
		AttributeID: attributeID,
		Required:    required,
		SortOrder:   sortOrder,
		Filterable:  filterable,
		Searchable:  searchable,
		Enabled:     enabled,
		CreatedAt:   createdAt,
		ModifiedAt:  modifiedAt,
	}
}

// Update modifies category attribute data with validation
func (ca *CategoryAttribute) Update(
	required bool,
	sortOrder int,
	filterable *bool,
	searchable *bool,
	enabled bool,
) error {
	if sortOrder < 0 {
		return errors.New("sortOrder cannot be negative")
	}

	ca.Required = required
	ca.SortOrder = sortOrder
	ca.Filterable = filterable
	ca.Searchable = searchable
	ca.Enabled = enabled
	ca.ModifiedAt = time.Now().UTC()

	return nil
}

func validateCategoryAttributeData(categoryID string, attributeID string, sortOrder int) error {
	if categoryID == "" {
		return errors.New("categoryID is required")
	}

	if attributeID == "" {
		return errors.New("attributeID is required")
	}

	if sortOrder < 0 {
		return errors.New("sortOrder cannot be negative")
	}

	return nil
}
