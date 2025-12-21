package categoryattribute

import (
	"context"

	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
)

type ListQuery struct {
	CategoryID string
	Page       int
	Size       int
	Enabled    *bool
	Filterable *bool
	Sort       string
	Order      string
}

type Repository interface {
	Insert(ctx context.Context, ca *CategoryAttribute) error

	FindByID(ctx context.Context, id string) (*CategoryAttribute, error)

	FindByCategoryAndAttribute(ctx context.Context, categoryID, attributeID string) (*CategoryAttribute, error)

	FindList(ctx context.Context, query ListQuery) (*commonsmongo.PageResult[CategoryAttribute], error)

	Update(ctx context.Context, ca *CategoryAttribute) (*CategoryAttribute, error)

	Delete(ctx context.Context, id string) error
}
