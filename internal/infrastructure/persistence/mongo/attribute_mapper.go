package mongo

import (
	"github.com/Sokol111/ecommerce-attribute-service/internal/domain/attribute"
)

type attributeMapper struct{}

func newAttributeMapper() *attributeMapper {
	return &attributeMapper{}
}

func (m *attributeMapper) ToEntity(a *attribute.Attribute) *attributeEntity {
	return &attributeEntity{
		ID:                a.ID,
		Version:           a.Version,
		Name:              a.Name,
		Slug:              a.Slug,
		Type:              string(a.Type),
		Unit:              a.Unit,
		DefaultFilterable: a.DefaultFilterable,
		DefaultSearchable: a.DefaultSearchable,
		SortOrder:         a.SortOrder,
		Enabled:           a.Enabled,
		CreatedAt:         a.CreatedAt,
		ModifiedAt:        a.ModifiedAt,
	}
}

func (m *attributeMapper) ToDomain(e *attributeEntity) *attribute.Attribute {
	return attribute.Reconstruct(
		e.ID,
		e.Version,
		e.Name,
		e.Slug,
		attribute.AttributeType(e.Type),
		e.Unit,
		e.DefaultFilterable,
		e.DefaultSearchable,
		e.SortOrder,
		e.Enabled,
		e.CreatedAt.UTC(),
		e.ModifiedAt.UTC(),
	)
}

func (m *attributeMapper) GetID(e *attributeEntity) string {
	return e.ID
}

func (m *attributeMapper) GetVersion(e *attributeEntity) int {
	return e.Version
}

func (m *attributeMapper) SetVersion(e *attributeEntity, version int) {
	e.Version = version
}
