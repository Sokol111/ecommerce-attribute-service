package attribute

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// AttributeType represents the type of attribute
type AttributeType string

const (
	AttributeTypeSelect      AttributeType = "select"
	AttributeTypeMultiselect AttributeType = "multiselect"
	AttributeTypeRange       AttributeType = "range"
	AttributeTypeBoolean     AttributeType = "boolean"
	AttributeTypeText        AttributeType = "text"
)

// Attribute - domain aggregate root
type Attribute struct {
	ID                string
	Version           int
	Name              string
	Slug              string
	Type              AttributeType
	Unit              *string
	DefaultFilterable bool
	DefaultSearchable bool
	SortOrder         int
	Enabled           bool
	CreatedAt         time.Time
	ModifiedAt        time.Time
}

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// NewAttribute creates a new attribute with validation.
// If id is empty, a new UUID will be generated.
func NewAttribute(
	id string,
	name string,
	slug string,
	attrType AttributeType,
	unit *string,
	defaultFilterable bool,
	defaultSearchable bool,
	sortOrder int,
	enabled bool,
) (*Attribute, error) {
	if err := validateAttributeData(name, slug, attrType, sortOrder); err != nil {
		return nil, err
	}

	if id == "" {
		id = uuid.New().String()
	}

	now := time.Now().UTC()
	return &Attribute{
		ID:                id,
		Version:           1,
		Name:              name,
		Slug:              slug,
		Type:              attrType,
		Unit:              unit,
		DefaultFilterable: defaultFilterable,
		DefaultSearchable: defaultSearchable,
		SortOrder:         sortOrder,
		Enabled:           enabled,
		CreatedAt:         now,
		ModifiedAt:        now,
	}, nil
}

// Reconstruct rebuilds an attribute from persistence (no validation)
func Reconstruct(
	id string,
	version int,
	name string,
	slug string,
	attrType AttributeType,
	unit *string,
	defaultFilterable bool,
	defaultSearchable bool,
	sortOrder int,
	enabled bool,
	createdAt time.Time,
	modifiedAt time.Time,
) *Attribute {
	return &Attribute{
		ID:                id,
		Version:           version,
		Name:              name,
		Slug:              slug,
		Type:              attrType,
		Unit:              unit,
		DefaultFilterable: defaultFilterable,
		DefaultSearchable: defaultSearchable,
		SortOrder:         sortOrder,
		Enabled:           enabled,
		CreatedAt:         createdAt,
		ModifiedAt:        modifiedAt,
	}
}

// Update modifies attribute data with validation
func (a *Attribute) Update(
	name string,
	slug string,
	attrType AttributeType,
	unit *string,
	defaultFilterable bool,
	defaultSearchable bool,
	sortOrder int,
	enabled bool,
) error {
	if err := validateAttributeData(name, slug, attrType, sortOrder); err != nil {
		return err
	}

	a.Name = name
	a.Slug = slug
	a.Type = attrType
	a.Unit = unit
	a.DefaultFilterable = defaultFilterable
	a.DefaultSearchable = defaultSearchable
	a.SortOrder = sortOrder
	a.Enabled = enabled
	a.ModifiedAt = time.Now().UTC()

	return nil
}

// validateAttributeData validates business rules
func validateAttributeData(name string, slug string, attrType AttributeType, sortOrder int) error {
	if name == "" {
		return errors.New("name is required")
	}

	if len(name) > 100 {
		return errors.New("name is too long (max 100 characters)")
	}

	if slug == "" {
		return errors.New("slug is required")
	}

	if len(slug) > 50 {
		return errors.New("slug is too long (max 50 characters)")
	}

	if !slugRegex.MatchString(slug) {
		return errors.New("slug must contain only lowercase letters, numbers, and hyphens")
	}

	if !isValidAttributeType(attrType) {
		return errors.New("invalid attribute type")
	}

	if sortOrder < 0 {
		return errors.New("sortOrder cannot be negative")
	}

	return nil
}

func isValidAttributeType(t AttributeType) bool {
	switch t {
	case AttributeTypeSelect, AttributeTypeMultiselect, AttributeTypeRange, AttributeTypeBoolean, AttributeTypeText:
		return true
	}
	return false
}
