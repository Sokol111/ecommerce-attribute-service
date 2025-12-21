package categoryattribute

import "errors"

var (
	ErrAlreadyAssigned = errors.New("attribute is already assigned to this category")
)
