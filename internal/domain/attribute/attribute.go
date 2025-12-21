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

// Option represents an attribute option (embedded in Attribute)
type Option struct {
	Value     string
	Slug      string
	ColorCode *string
	SortOrder int
	Enabled   bool
}

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
	Options           []Option
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
	options []Option,
) (*Attribute, error) {
	if err := validateAttributeData(name, slug, attrType, sortOrder); err != nil {
		return nil, err
	}

	if err := validateOptions(options); err != nil {
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
		Options:           options,
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
	options []Option,
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
		Options:           options,
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
	options []Option,
) error {
	if err := validateAttributeData(name, slug, attrType, sortOrder); err != nil {
		return err
	}

	if err := validateOptions(options); err != nil {
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
	a.Options = options
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

// validateOptions validates option data
func validateOptions(options []Option) error {
	if len(options) == 0 {
		return nil
	}

	slugs := make(map[string]bool)
	for _, opt := range options {
		if opt.Value == "" {
			return errors.New("option value is required")
		}
		if len(opt.Value) > 100 {
			return errors.New("option value is too long (max 100 characters)")
		}
		if opt.Slug == "" {
			return errors.New("option slug is required")
		}
		if len(opt.Slug) > 50 {
			return errors.New("option slug is too long (max 50 characters)")
		}
		if !slugRegex.MatchString(opt.Slug) {
			return errors.New("option slug must contain only lowercase letters, numbers, and hyphens")
		}
		if slugs[opt.Slug] {
			return errors.New("duplicate option slug: " + opt.Slug)
		}
		slugs[opt.Slug] = true
		if opt.SortOrder < 0 {
			return errors.New("option sortOrder cannot be negative")
		}
	}
	return nil
}
