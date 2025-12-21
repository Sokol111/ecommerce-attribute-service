package mongo

import (
	"context"

	"github.com/Sokol111/ecommerce-attribute-service/internal/domain/categoryattribute"
	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type categoryAttributeRepository struct {
	*commonsmongo.GenericRepository[categoryattribute.CategoryAttribute, categoryAttributeEntity]
}

func newCategoryAttributeRepository(mongoClient commonsmongo.Mongo, mapper *categoryAttributeMapper) (categoryattribute.Repository, error) {
	collection := mongoClient.GetCollection("category_attribute")

	genericRepo, err := commonsmongo.NewGenericRepository(
		collection,
		mapper,
	)
	if err != nil {
		return nil, err
	}

	return &categoryAttributeRepository{
		GenericRepository: genericRepo,
	}, nil
}

func (r *categoryAttributeRepository) FindByCategoryAndAttribute(ctx context.Context, categoryID, attributeID string) (*categoryattribute.CategoryAttribute, error) {
	filter := bson.D{
		{Key: "categoryId", Value: categoryID},
		{Key: "attributeId", Value: attributeID},
	}

	opts := commonsmongo.QueryOptions{
		Filter: filter,
		Page:   1,
		Size:   1,
	}

	result, err := r.FindWithOptions(ctx, opts)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, persistence.ErrEntityNotFound
	}

	return result.Items[0], nil
}

func (r *categoryAttributeRepository) FindList(ctx context.Context, query categoryattribute.ListQuery) (*commonsmongo.PageResult[categoryattribute.CategoryAttribute], error) {
	filter := bson.D{{Key: "categoryId", Value: query.CategoryID}}

	if query.Enabled != nil {
		filter = append(filter, bson.E{Key: "enabled", Value: *query.Enabled})
	}
	if query.Filterable != nil {
		filter = append(filter, bson.E{Key: "filterable", Value: *query.Filterable})
	}

	var sortBson bson.D
	if query.Sort != "" {
		sortOrder := 1 // asc
		if query.Order == "desc" {
			sortOrder = -1
		}
		sortBson = bson.D{{Key: query.Sort, Value: sortOrder}}
	}

	opts := commonsmongo.QueryOptions{
		Filter: filter,
		Page:   query.Page,
		Size:   query.Size,
		Sort:   sortBson,
	}

	return r.FindWithOptions(ctx, opts)
}

// Override Insert to handle duplicate assignment error
func (r *categoryAttributeRepository) Insert(ctx context.Context, ca *categoryattribute.CategoryAttribute) error {
	err := r.GenericRepository.Insert(ctx, ca)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return categoryattribute.ErrAlreadyAssigned
		}
		return err
	}
	return nil
}

// Override Update to handle duplicate assignment error
func (r *categoryAttributeRepository) Update(ctx context.Context, ca *categoryattribute.CategoryAttribute) (*categoryattribute.CategoryAttribute, error) {
	result, err := r.GenericRepository.Update(ctx, ca)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, categoryattribute.ErrAlreadyAssigned
		}
		return nil, err
	}
	return result, nil
}
