package application

import (
	"github.com/Sokol111/ecommerce-attribute-service/internal/application/command"
	"github.com/Sokol111/ecommerce-attribute-service/internal/application/query"
	"go.uber.org/fx"
)

// Module provides application layer dependencies
func Module() fx.Option {
	return fx.Options(
		// Command handlers
		fx.Provide(
			command.NewCreateAttributeHandler,
			command.NewUpdateAttributeHandler,
			command.NewAssignAttributeToCategoryHandler,
			command.NewUpdateCategoryAttributeHandler,
			command.NewUnassignAttributeFromCategoryHandler,
		),
		// Query handlers
		fx.Provide(
			query.NewGetAttributeByIDHandler,
			query.NewGetAttributeListHandler,
			query.NewGetCategoryAttributeListHandler,
		),
	)
}
