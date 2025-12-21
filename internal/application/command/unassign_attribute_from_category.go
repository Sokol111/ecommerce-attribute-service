package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sokol111/ecommerce-attribute-service/internal/domain/categoryattribute"
	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
)

type UnassignAttributeFromCategoryCommand struct {
	ID         string
	CategoryID string // for validation
}

type UnassignAttributeFromCategoryCommandHandler interface {
	Handle(ctx context.Context, cmd UnassignAttributeFromCategoryCommand) error
}

type unassignAttributeFromCategoryHandler struct {
	repo categoryattribute.Repository
}

func NewUnassignAttributeFromCategoryHandler(repo categoryattribute.Repository) UnassignAttributeFromCategoryCommandHandler {
	return &unassignAttributeFromCategoryHandler{
		repo: repo,
	}
}

func (h *unassignAttributeFromCategoryHandler) Handle(ctx context.Context, cmd UnassignAttributeFromCategoryCommand) error {
	// Verify the assignment exists and belongs to the category
	ca, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		if !errors.Is(err, persistence.ErrEntityNotFound) {
			return fmt.Errorf("failed to get category attribute: %w", err)
		}
		return err
	}

	if ca.CategoryID != cmd.CategoryID {
		return persistence.ErrEntityNotFound
	}

	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete category attribute: %w", err)
	}

	return nil
}
