package command

import (
	"context"
	"fmt"

	"github.com/Sokol111/ecommerce-attribute-service/internal/domain/attribute"
	"github.com/Sokol111/ecommerce-attribute-service/internal/domain/categoryattribute"
	"github.com/Sokol111/ecommerce-commons/pkg/persistence"
)

type AssignAttributeToCategoryCommand struct {
	ID          *string
	CategoryID  string
	AttributeID string
	Required    bool
	SortOrder   int
	Filterable  *bool
	Searchable  *bool
	Enabled     bool
}

type AssignAttributeToCategoryCommandHandler interface {
	Handle(ctx context.Context, cmd AssignAttributeToCategoryCommand) (*categoryattribute.CategoryAttribute, error)
}

type assignAttributeToCategoryHandler struct {
	caRepo   categoryattribute.Repository
	attrRepo attribute.Repository
}

func NewAssignAttributeToCategoryHandler(
	caRepo categoryattribute.Repository,
	attrRepo attribute.Repository,
) AssignAttributeToCategoryCommandHandler {
	return &assignAttributeToCategoryHandler{
		caRepo:   caRepo,
		attrRepo: attrRepo,
	}
}

func (h *assignAttributeToCategoryHandler) Handle(ctx context.Context, cmd AssignAttributeToCategoryCommand) (*categoryattribute.CategoryAttribute, error) {
	// Verify attribute exists
	exists, err := h.attrRepo.Exists(ctx, cmd.AttributeID)
	if err != nil {
		return nil, fmt.Errorf("failed to check attribute existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("attribute not found: %w", persistence.ErrEntityNotFound)
	}

	var id string
	if cmd.ID != nil {
		id = *cmd.ID
	}

	ca, err := categoryattribute.NewCategoryAttribute(
		id,
		cmd.CategoryID,
		cmd.AttributeID,
		cmd.Required,
		cmd.SortOrder,
		cmd.Filterable,
		cmd.Searchable,
		cmd.Enabled,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create category attribute: %w", err)
	}

	if err := h.caRepo.Insert(ctx, ca); err != nil {
		return nil, fmt.Errorf("failed to insert category attribute: %w", err)
	}

	return ca, nil
}
