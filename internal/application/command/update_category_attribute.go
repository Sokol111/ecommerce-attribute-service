package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sokol111/ecommerce-attribute-service/internal/domain/categoryattribute"
	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
)

type UpdateCategoryAttributeCommand struct {
	ID         string
	CategoryID string // for validation
	Version    int
	Required   bool
	SortOrder  int
	Filterable *bool
	Searchable *bool
	Enabled    bool
}

type UpdateCategoryAttributeCommandHandler interface {
	Handle(ctx context.Context, cmd UpdateCategoryAttributeCommand) (*categoryattribute.CategoryAttribute, error)
}

type updateCategoryAttributeHandler struct {
	repo categoryattribute.Repository
}

func NewUpdateCategoryAttributeHandler(repo categoryattribute.Repository) UpdateCategoryAttributeCommandHandler {
	return &updateCategoryAttributeHandler{
		repo: repo,
	}
}

func (h *updateCategoryAttributeHandler) Handle(ctx context.Context, cmd UpdateCategoryAttributeCommand) (*categoryattribute.CategoryAttribute, error) {
	ca, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if !errors.Is(err, persistence.ErrEntityNotFound) {
			return nil, fmt.Errorf("failed to get category attribute: %w", err)
		}
		return nil, err
	}

	// Verify the assignment belongs to the specified category
	if ca.CategoryID != cmd.CategoryID {
		return nil, persistence.ErrEntityNotFound
	}

	if ca.Version != cmd.Version {
		return nil, persistence.ErrOptimisticLocking
	}

	if err := ca.Update(
		cmd.Required,
		cmd.SortOrder,
		cmd.Filterable,
		cmd.Searchable,
		cmd.Enabled,
	); err != nil {
		return nil, fmt.Errorf("failed to update category attribute: %w", err)
	}

	updated, err := h.repo.Update(ctx, ca)
	if err != nil {
		if !errors.Is(err, persistence.ErrOptimisticLocking) {
			return nil, fmt.Errorf("failed to update category attribute: %w", err)
		}
		return nil, err
	}

	return updated, nil
}
