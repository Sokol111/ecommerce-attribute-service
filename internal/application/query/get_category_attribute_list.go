package query

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-attribute-service/internal/domain/categoryattribute"
)

type GetCategoryAttributeListQuery struct {
	CategoryID string
	Page       int
	Size       int
	Enabled    *bool
	Filterable *bool
	Sort       string
	Order      string
}

type ListCategoryAttributesResult struct {
	Items []*categoryattribute.CategoryAttribute
	Page  int
	Size  int
	Total int64
}

type GetCategoryAttributeListQueryHandler interface {
	Handle(ctx context.Context, query GetCategoryAttributeListQuery) (*ListCategoryAttributesResult, error)
}

type getCategoryAttributeListHandler struct {
	repo categoryattribute.Repository
}

func NewGetCategoryAttributeListHandler(repo categoryattribute.Repository) GetCategoryAttributeListQueryHandler {
	return &getCategoryAttributeListHandler{repo: repo}
}

func (h *getCategoryAttributeListHandler) Handle(ctx context.Context, query GetCategoryAttributeListQuery) (*ListCategoryAttributesResult, error) {
	listQuery := categoryattribute.ListQuery{
		CategoryID: query.CategoryID,
		Page:       query.Page,
		Size:       query.Size,
		Enabled:    query.Enabled,
		Filterable: query.Filterable,
		Sort:       query.Sort,
		Order:      query.Order,
	}

	result, err := h.repo.FindList(ctx, listQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get category attributes list: %w", err)
	}

	return &ListCategoryAttributesResult{
		Items: result.Items,
		Page:  result.Page,
		Size:  result.Size,
		Total: result.Total,
	}, nil
}
