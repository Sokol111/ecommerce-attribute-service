package mongo

import (
	"github.com/Sokol111/ecommerce-attribute-service/internal/domain/categoryattribute"
)

type categoryAttributeMapper struct{}

func newCategoryAttributeMapper() *categoryAttributeMapper {
	return &categoryAttributeMapper{}
}

func (m *categoryAttributeMapper) ToEntity(ca *categoryattribute.CategoryAttribute) *categoryAttributeEntity {
	return &categoryAttributeEntity{
		ID:          ca.ID,
		Version:     ca.Version,
		CategoryID:  ca.CategoryID,
		AttributeID: ca.AttributeID,
		Required:    ca.Required,
		SortOrder:   ca.SortOrder,
		Filterable:  ca.Filterable,
		Searchable:  ca.Searchable,
		Enabled:     ca.Enabled,
		CreatedAt:   ca.CreatedAt,
		ModifiedAt:  ca.ModifiedAt,
	}
}

func (m *categoryAttributeMapper) ToDomain(e *categoryAttributeEntity) *categoryattribute.CategoryAttribute {
	return categoryattribute.Reconstruct(
		e.ID,
		e.Version,
		e.CategoryID,
		e.AttributeID,
		e.Required,
		e.SortOrder,
		e.Filterable,
		e.Searchable,
		e.Enabled,
		e.CreatedAt.UTC(),
		e.ModifiedAt.UTC(),
	)
}

func (m *categoryAttributeMapper) GetID(e *categoryAttributeEntity) string {
	return e.ID
}

func (m *categoryAttributeMapper) GetVersion(e *categoryAttributeEntity) int {
	return e.Version
}

func (m *categoryAttributeMapper) SetVersion(e *categoryAttributeEntity, version int) {
	e.Version = version
}
